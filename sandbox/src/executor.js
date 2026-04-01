const ivm = require("isolated-vm");

// メモリ制限: 128MB
const MEMORY_LIMIT_MB = 128;

// 出力制限（合計・件数・アイテムごと）
const MAX_OUTPUT_ITEMS = 100;
const MAX_OUTPUT_CHARS = 64 * 1024; // 64KB
const MAX_PER_ITEM_CHARS = 8 * 1024; // 8KB
const TRUNC_MSG = "... (truncated)";

function truncateString(s, max) {
  try {
    if (typeof s !== "string") s = String(s);
    if (s.length <= max) return s;
    return s.slice(0, max) + TRUNC_MSG;
  } catch {
    return "[Unserializable]";
  }
}

function safeSerializeShallow(v) {
  try {
    if (v === null) return "null";
    if (v === undefined) return "undefined";
    const t = typeof v;
    if (t === "string") return v.length > 200 ? v.slice(0, 200) + "..." : v;
    if (t === "number" || t === "boolean" || t === "bigint") return String(v);
    if (t === "function") return "[Function]";
    if (Array.isArray(v)) return `[Array(${v.length})]`;
    if (t === "object") {
      try {
        const keys = Object.keys(v);
        const entries = keys
          .slice(0, 5)
          .map((k) => `${k}:${safeSerializeShallow(v[k])}`)
          .join(", ");
        return `{${entries}${keys.length > 5 ? ",..." : ""}}`;
      } catch {
        return "[Object]";
      }
    }
    return String(v);
  } catch {
    return "[Unserializable]";
  }
}

function safeSerialize(value) {
  try {
    if (value === null) return "null";
    if (value === undefined) return "undefined";
    const t = typeof value;
    if (t === "string") return truncateString(value, MAX_PER_ITEM_CHARS);
    if (t === "number" || t === "boolean" || t === "bigint") return String(value);
    if (t === "function") return "[Function]";
    if (Array.isArray(value)) {
      const preview = value
        .slice(0, 5)
        .map((v) => safeSerializeShallow(v))
        .join(", ");
      return truncateString(`[Array(${value.length})] ${preview}`, MAX_PER_ITEM_CHARS);
    }
    if (t === "object") {
      const keys = Object.keys(value || {});
      const entries = keys
        .slice(0, 5)
        .map((k) => `${k}:${safeSerializeShallow(value[k])}`)
        .join(", ");
      return truncateString(`{${entries}${keys.length > 5 ? ",..." : ""}}`, MAX_PER_ITEM_CHARS);
    }
    return truncateString(String(value), MAX_PER_ITEM_CHARS);
  } catch {
    return "[Unserializable]";
  }
}

/**
 * 隔離された環境でJavaScriptコードを実行する
 * @param {string} code - 実行するコード
 * @param {number} timeout - タイムアウト（ミリ秒）
 * @returns {Promise<{success: boolean, output: Array, error?: string}>}
 */
async function executeCode(code, timeout = 10000) {
  const output = [];
  let outputChars = 0;
  let isolate = null;
  let context = null;

  console.log(`[Sandbox] Creating new Isolate for execution...`);

  try {
    isolate = new ivm.Isolate({ memoryLimit: MEMORY_LIMIT_MB });
    console.log(`[Sandbox] Isolate created (memory limit: ${MEMORY_LIMIT_MB}MB)`);

    context = await isolate.createContext();

    const jail = context.global;

    await jail.set("global", jail.derefInto());

    // console.log（軽量化・切り捨てを行う）
    function pushOutput(item) {
      try {
        if (output.length >= MAX_OUTPUT_ITEMS) return;
        const text = item && item.text ? String(item.text) : "";
        const truncated = truncateString(text, MAX_PER_ITEM_CHARS);
        const newChars = outputChars + truncated.length;
        if (newChars > MAX_OUTPUT_CHARS) {
          const remaining = MAX_OUTPUT_CHARS - outputChars;
          if (remaining <= 0) return;
          item.text = truncated.slice(0, remaining) + TRUNC_MSG;
          output.push(item);
          outputChars = MAX_OUTPUT_CHARS;
          return;
        }
        item.text = truncated;
        output.push(item);
        outputChars = newChars;
      } catch {
        // ここで例外を投げるとホストが影響を受けるため無視する
      }
    }

    const logCallback = new ivm.Callback((...args) => {
      try {
        const parts = args.map((arg) => {
          try {
            if (arg && typeof arg === "object") return safeSerialize(arg);
          } catch {}
          try {
            return String(arg);
          } catch {
            return "[Unserializable]";
          }
        });
        pushOutput({ type: "log", text: parts.join(" ") });
      } catch {
        pushOutput({ type: "error", text: "log callback error" });
      }
    });

    const errorCallback = new ivm.Callback((...args) => {
      try {
        const parts = args.map((arg) => {
          try {
            if (arg && typeof arg === "object") return safeSerialize(arg);
          } catch {}
          try {
            return String(arg);
          } catch {
            return "[Unserializable]";
          }
        });
        pushOutput({ type: "error", text: parts.join(" ") });
      } catch {
        // ignore
      }
    });

    const warnCallback = new ivm.Callback((...args) => {
      try {
        const parts = args.map((arg) => {
          try {
            if (arg && typeof arg === "object") return safeSerialize(arg);
          } catch {}
          try {
            return String(arg);
          } catch {
            return "[Unserializable]";
          }
        });
        pushOutput({ type: "warn", text: parts.join(" ") });
      } catch {}
    });

    const infoCallback = new ivm.Callback((...args) => {
      try {
        const parts = args.map((arg) => {
          try {
            if (arg && typeof arg === "object") return safeSerialize(arg);
          } catch {}
          try {
            return String(arg);
          } catch {
            return "[Unserializable]";
          }
        });
        pushOutput({ type: "info", text: parts.join(" ") });
      } catch {}
    });

    // consoleオブジェクトを設定
    await jail.set("_log", logCallback);
    await jail.set("_error", errorCallback);
    await jail.set("_warn", warnCallback);
    await jail.set("_info", infoCallback);

    // consoleオブジェクトを作成するブートストラップコード
    const bootstrap = `
      const console = {
        log: (...args) => _log(...args),
        error: (...args) => _error(...args),
        warn: (...args) => _warn(...args),
        info: (...args) => _info(...args),
        debug: (...args) => _log(...args),
        trace: (...args) => _log(...args),
      };
      
      // 危険なグローバルを削除/無効化
      const _setTimeout = undefined;
      const _setInterval = undefined;
      const _setImmediate = undefined;
      const _fetch = undefined;
      const _XMLHttpRequest = undefined;
      const _WebSocket = undefined;
    `;

    await context.eval(bootstrap);

    const script = await isolate.compileScript(code);

    const result = await script.run(context, { timeout });

    // 結果があれば出力に追加（安全にシリアライズ）
    if (result !== undefined) {
      const resultText = safeSerialize(result);
      pushOutput({ type: "result", text: `=> ${resultText}` });
    }

    return {
      success: true,
      output,
    };
  } catch (error) {
    let errorMessage = error.message || "Unknown error";

    // タイムアウトエラーの判定
    if (errorMessage.includes("Script execution timed out")) {
      errorMessage = `Execution timed out (limit: ${timeout}ms)`;
    }

    // メモリ制限エラーの判定
    if (errorMessage.includes("Isolate was disposed")) {
      errorMessage = `Memory limit exceeded (limit: ${MEMORY_LIMIT_MB}MB)`;
    }

    pushOutput({ type: "error", text: `Error: ${errorMessage}` });

    return {
      success: false,
      output,
      error: errorMessage,
    };
  } finally {
    if (context) {
      context.release();
      console.log(`[Sandbox] Context released`);
    }
    if (isolate) {
      isolate.dispose();
      console.log(`[Sandbox] Isolate disposed - cleanup complete`);
    }
  }
}

module.exports = { executeCode };
