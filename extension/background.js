// VocabWeb Extension - Background Service Worker

const API_BASE_URL_KEY = 'vocabweb_api_url';
const API_TOKEN_KEY = 'vocabweb_token';
const DEFAULT_API_URL = 'https://api.vocabweb.example.com';

// 安装时初始化
chrome.runtime.onInstalled.addListener(() => {
  console.log('VocabWeb Extension installed');
  
  // 创建右键菜单
  chrome.contextMenus.create({
    id: 'vocabweb-add-word',
    title: '收藏到 VocabWeb',
    contexts: ['selection']
  });
});

// 监听右键菜单点击
chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === 'vocabweb-add-word' && info.selectionText) {
    const word = info.selectionText.trim();
    addWordToVocabWeb(word, tab.id);
  }
});

// 监听来自 content script 的消息
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'addWord') {
    addWordToVocabWeb(request.word, sender.tab.id)
      .then(result => sendResponse({ success: true, data: result }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true; // 保持消息通道开启
  }
  
  if (request.action === 'getConfig') {
    getConfig()
      .then(config => sendResponse({ success: true, config }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true;
  }
});

// 获取配置
async function getConfig() {
  const result = await chrome.storage.local.get([API_BASE_URL_KEY, API_TOKEN_KEY]);
  return {
    apiUrl: result[API_BASE_URL_KEY] || DEFAULT_API_URL,
    token: result[API_TOKEN_KEY] || ''
  };
}

// 添加单词到 VocabWeb
async function addWordToVocabWeb(word, tabId) {
  const config = await getConfig();
  
  if (!config.token) {
    throw new Error('请先在插件设置中配置 API Token');
  }
  
  // 清理单词（只保留字母）
  const cleanWord = word.replace(/[^a-zA-Z]/g, '').toLowerCase();
  if (!cleanWord) {
    throw new Error('无效的单词');
  }
  
  // 调用后端 API
  const response = await fetch(`${config.apiUrl}/api/words`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${config.token}`
    },
    body: JSON.stringify({ word: cleanWord })
  });
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(`添加失败: ${error}`);
  }
  
  const result = await response.json();
  
  // 通知 content script 显示成功消息
  if (tabId) {
    chrome.tabs.sendMessage(tabId, {
      action: 'showSuccess',
      word: cleanWord
    });
  }
  
  return result;
}
