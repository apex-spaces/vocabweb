// VocabWeb Extension - Popup Script

const API_BASE_URL_KEY = 'vocabweb_api_url';
const API_TOKEN_KEY = 'vocabweb_token';
const DEFAULT_API_URL = 'https://api.vocabweb.example.com';

// DOM 元素
const loginSection = document.getElementById('login-section');
const dashboardSection = document.getElementById('dashboard-section');
const apiUrlInput = document.getElementById('api-url');
const apiTokenInput = document.getElementById('api-token');
const saveConfigBtn = document.getElementById('save-config-btn');
const logoutBtn = document.getElementById('logout-btn');
const todayCountEl = document.getElementById('today-count');
const totalCountEl = document.getElementById('total-count');
const recentListEl = document.getElementById('recent-list');

// 初始化
async function init() {
  const config = await loadConfig();
  
  if (config.token) {
    showDashboard();
    loadDashboardData(config);
  } else {
    showLogin();
    apiUrlInput.value = config.apiUrl || DEFAULT_API_URL;
  }
}

// 加载配置
async function loadConfig() {
  const result = await chrome.storage.local.get([API_BASE_URL_KEY, API_TOKEN_KEY]);
  return {
    apiUrl: result[API_BASE_URL_KEY] || DEFAULT_API_URL,
    token: result[API_TOKEN_KEY] || ''
  };
}

// 保存配置
async function saveConfig(apiUrl, token) {
  await chrome.storage.local.set({
    [API_BASE_URL_KEY]: apiUrl,
    [API_TOKEN_KEY]: token
  });
}

// 显示登录界面
function showLogin() {
  loginSection.classList.remove('hidden');
  dashboardSection.classList.add('hidden');
}

// 显示仪表盘
function showDashboard() {
  loginSection.classList.add('hidden');
  dashboardSection.classList.remove('hidden');
}

// 加载仪表盘数据
async function loadDashboardData(config) {
  try {
    // 获取统计数据
    const statsResponse = await fetch(`${config.apiUrl}/api/stats`, {
      headers: {
        'Authorization': `Bearer ${config.token}`
      }
    });
    
    if (statsResponse.ok) {
      const stats = await statsResponse.json();
      todayCountEl.textContent = stats.today || 0;
      totalCountEl.textContent = stats.total || 0;
    }
    
    // 获取最近单词
    const wordsResponse = await fetch(`${config.apiUrl}/api/words?limit=5`, {
      headers: {
        'Authorization': `Bearer ${config.token}`
      }
    });
    
    if (wordsResponse.ok) {
      const words = await wordsResponse.json();
      renderRecentWords(words);
    }
  } catch (error) {
    recentListEl.innerHTML = '<div class="loading">加载失败</div>';
  }
}

// 渲染最近单词列表
function renderRecentWords(words) {
  if (!words || words.length === 0) {
    recentListEl.innerHTML = '<div class="loading">暂无收藏</div>';
    return;
  }
  
  recentListEl.innerHTML = words.map(word => `
    <div class="word-item">
      <div class="word-name">${escapeHtml(word.word)}</div>
      <div class="word-time">${formatTime(word.created_at)}</div>
    </div>
  `).join('');
}

// 格式化时间
function formatTime(timestamp) {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now - date;
  
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`;
  return date.toLocaleDateString('zh-CN');
}

// HTML 转义
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// 事件监听
saveConfigBtn.addEventListener('click', async () => {
  const apiUrl = apiUrlInput.value.trim();
  const token = apiTokenInput.value.trim();
  
  if (!apiUrl || !token) {
    alert('请填写完整信息');
    return;
  }
  
  saveConfigBtn.textContent = '保存中...';
  saveConfigBtn.disabled = true;
  
  try {
    await saveConfig(apiUrl, token);
    alert('配置已保存');
    location.reload();
  } catch (error) {
    alert('保存失败: ' + error.message);
    saveConfigBtn.textContent = '保存配置';
    saveConfigBtn.disabled = false;
  }
});

logoutBtn.addEventListener('click', async () => {
  if (confirm('确定要退出登录吗？')) {
    await chrome.storage.local.remove([API_TOKEN_KEY]);
    location.reload();
  }
});

// 启动
init();
