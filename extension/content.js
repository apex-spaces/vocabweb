// VocabWeb Extension - Content Script
// 监听页面文本选中，显示浮动气泡

(function() {
  'use strict';
  
  let bubble = null;
  let currentSelection = '';
  
  // 初始化
  function init() {
    document.addEventListener('mouseup', handleTextSelection);
    document.addEventListener('mousedown', hideBubble);
    
    // 监听来自 background 的消息
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
      if (request.action === 'showSuccess') {
        showSuccessMessage(request.word);
      }
    });
  }
  
  // 处理文本选中事件
  function handleTextSelection(e) {
    // 延迟执行，确保选中完成
    setTimeout(() => {
      const selection = window.getSelection();
      const text = selection.toString().trim();
      
      // 检查是否选中了文本
      if (!text) {
        hideBubble();
        return;
      }
      
      // 检查是否是英文单词（允许连字符）
      const word = text.replace(/[^a-zA-Z-]/g, '');
      if (!word || word.length < 2) {
        hideBubble();
        return;
      }
      
      currentSelection = word;
      showBubble(word, e.pageX, e.pageY);
    }, 10);
  }
  
  // 显示浮动气泡
  function showBubble(word, x, y) {
    hideBubble();
    
    bubble = document.createElement('div');
    bubble.className = 'vocabweb-bubble';
    bubble.innerHTML = `
      <div class="vocabweb-bubble-word">${escapeHtml(word)}</div>
      <div class="vocabweb-bubble-loading">正在查询释义...</div>
      <button class="vocabweb-bubble-btn">收藏到 VocabWeb</button>
    `;
    
    document.body.appendChild(bubble);
    
    // 定位气泡
    positionBubble(bubble, x, y);
    
    // 绑定收藏按钮
    const btn = bubble.querySelector('.vocabweb-bubble-btn');
    btn.addEventListener('click', () => handleAddWord(word));
    
    // 获取释义
    fetchDefinition(word);
  }
  
  // 定位气泡
  function positionBubble(bubble, x, y) {
    const rect = bubble.getBoundingClientRect();
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;
    
    let left = x + 10;
    let top = y + 10;
    
    // 防止超出右边界
    if (left + rect.width > viewportWidth) {
      left = x - rect.width - 10;
    }
    
    // 防止超出下边界
    if (top + rect.height > viewportHeight + window.scrollY) {
      top = y - rect.height - 10;
    }
    
    bubble.style.left = left + 'px';
    bubble.style.top = top + 'px';
  }
  
  // 隐藏气泡
  function hideBubble() {
    if (bubble && bubble.parentNode) {
      bubble.parentNode.removeChild(bubble);
      bubble = null;
    }
  }
  
  // 获取单词释义（使用免费词典API）
  async function fetchDefinition(word) {
    if (!bubble) return;
    
    const loadingEl = bubble.querySelector('.vocabweb-bubble-loading');
    
    try {
      const response = await fetch(`https://api.dictionaryapi.dev/api/v2/entries/en/${word}`);
      
      if (!response.ok) {
        loadingEl.textContent = '未找到释义';
        return;
      }
      
      const data = await response.json();
      const firstEntry = data[0];
      const meaning = firstEntry.meanings[0];
      const definition = meaning.definitions[0].definition;
      
      loadingEl.innerHTML = `<span class="vocabweb-bubble-def">${escapeHtml(definition)}</span>`;
    } catch (error) {
      loadingEl.textContent = '释义加载失败';
    }
  }
  
  // 处理添加单词
  async function handleAddWord(word) {
    if (!bubble) return;
    
    const btn = bubble.querySelector('.vocabweb-bubble-btn');
    btn.disabled = true;
    btn.textContent = '收藏中...';
    
    try {
      const response = await chrome.runtime.sendMessage({
        action: 'addWord',
        word: word
      });
      
      if (response.success) {
        btn.textContent = '✓ 已收藏';
        btn.classList.add('vocabweb-bubble-btn-success');
        setTimeout(hideBubble, 2000);
      } else {
        btn.textContent = '收藏失败';
        btn.classList.add('vocabweb-bubble-btn-error');
        alert(response.error || '收藏失败，请检查配置');
      }
    } catch (error) {
      btn.textContent = '收藏失败';
      btn.classList.add('vocabweb-bubble-btn-error');
      alert('收藏失败: ' + error.message);
    }
  }
  
  // 显示成功消息
  function showSuccessMessage(word) {
    const msg = document.createElement('div');
    msg.className = 'vocabweb-success-msg';
    msg.textContent = `✓ "${word}" 已收藏`;
    document.body.appendChild(msg);
    
    setTimeout(() => {
      msg.classList.add('vocabweb-success-msg-show');
    }, 10);
    
    setTimeout(() => {
      msg.classList.remove('vocabweb-success-msg-show');
      setTimeout(() => {
        if (msg.parentNode) {
          msg.parentNode.removeChild(msg);
        }
      }, 300);
    }, 3000);
  }
  
  // HTML 转义
  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
  
  // 启动
  init();
})();


