const ivm = require('isolated-vm');

// メモリ制限: 128MB
const MEMORY_LIMIT_MB = 128;

/**
 * 隔離された環境でJavaScriptコードを実行する
 * @param {string} code - 実行するコード
 * @param {number} timeout - タイムアウト（ミリ秒）
 * @returns {Promise<{success: boolean, output: Array, error?: string}>}
 */
async function executeCode(code, timeout = 10000) {
  const output = [];
  let isolate = null;
  let context = null;

  try {
    isolate = new ivm.Isolate({ memoryLimit: MEMORY_LIMIT_MB });
    
    context = await isolate.createContext();
    
    const jail = context.global;
    
    await jail.set('global', jail.derefInto());

    // console.log
    const logCallback = new ivm.Callback((...args) => {
      const message = args.map(arg => {
        if (typeof arg === 'object') {
          try {
            return JSON.stringify(arg);
          } catch {
            return String(arg);
          }
        }
        return String(arg);
      }).join(' ');
      output.push({ type: 'log', text: message });
    });

    const errorCallback = new ivm.Callback((...args) => {
      const message = args.map(arg => String(arg)).join(' ');
      output.push({ type: 'error', text: message });
    });

    const warnCallback = new ivm.Callback((...args) => {
      const message = args.map(arg => String(arg)).join(' ');
      output.push({ type: 'warn', text: message });
    });

    const infoCallback = new ivm.Callback((...args) => {
      const message = args.map(arg => String(arg)).join(' ');
      output.push({ type: 'info', text: message });
    });

    // consoleオブジェクトを設定
    await jail.set('_log', logCallback);
    await jail.set('_error', errorCallback);
    await jail.set('_warn', warnCallback);
    await jail.set('_info', infoCallback);

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

    // 結果があれば出力に追加
    if (result !== undefined) {
      let resultText;
      if (typeof result === 'object') {
        try {
          resultText = JSON.stringify(result);
        } catch {
          resultText = String(result);
        }
      } else {
        resultText = String(result);
      }
      output.push({ type: 'result', text: `=> ${resultText}` });
    }

    return {
      success: true,
      output
    };

  } catch (error) {
    let errorMessage = error.message || 'Unknown error';
    
    // タイムアウトエラーの判定
    if (errorMessage.includes('Script execution timed out')) {
      errorMessage = `Execution timed out (limit: ${timeout}ms)`;
    }
    
    // メモリ制限エラーの判定
    if (errorMessage.includes('Isolate was disposed')) {
      errorMessage = `Memory limit exceeded (limit: ${MEMORY_LIMIT_MB}MB)`;
    }

    output.push({ type: 'error', text: `Error: ${errorMessage}` });

    return {
      success: false,
      output,
      error: errorMessage
    };

  } finally {
    if (context) {
      context.release();
    }
    if (isolate) {
      isolate.dispose();
    }
  }
}

module.exports = { executeCode };
