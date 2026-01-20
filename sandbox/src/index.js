const express = require('express');
const { executeCode } = require('./executor');

const app = express();
app.use(express.json({ limit: '1mb' }));

// ヘルスチェック
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// コード実行エンドポイント
app.post('/execute', async (req, res) => {
  const { code, timeout = 10000 } = req.body;

  if (!code || typeof code !== 'string') {
    return res.status(400).json({
      success: false,
      error: 'Code is required and must be a string',
      output: []
    });
  }

  // タイムアウトは最大10秒
  const actualTimeout = Math.min(timeout, 10000);

  try {
    const result = await executeCode(code, actualTimeout);
    res.json(result);
  } catch (error) {
    res.status(500).json({
      success: false,
      error: error.message || 'Internal server error',
      output: []
    });
  }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`Sandbox service listening on port ${PORT}`);
});
