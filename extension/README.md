# VocabWeb Chrome Extension

VocabWeb 的 Chrome 浏览器划词插件，帮助你在浏览网页时快速收藏英语单词。

## 功能特性

- 🖱️ **划词收藏**：选中英文单词即可显示浮动气泡
- 📖 **即时释义**：自动获取单词释义（使用免费词典 API）
- ⭐ **一键收藏**：点击按钮即可将单词同步到 VocabWeb
- 🎨 **深色主题**：与现代网站风格一致的深色界面
- 📊 **统计面板**：查看今日收藏数和最近单词

## 安装方法

### 1. 下载扩展

克隆或下载本项目到本地：

```bash
git clone <repository-url>
cd vocabweb-repo/extension
```

### 2. 加载到 Chrome

1. 打开 Chrome 浏览器
2. 访问 `chrome://extensions/`
3. 开启右上角的「开发者模式」
4. 点击「加载已解压的扩展程序」
5. 选择 `extension` 文件夹

### 3. 配置 API Token

1. 点击浏览器工具栏中的 VocabWeb 图标
2. 输入你的 API 地址（默认：`https://api.vocabweb.example.com`）
3. 输入你的 API Token（从 VocabWeb 后台获取）
4. 点击「保存配置」

## 使用方法

### 划词收藏

1. 在任意网页上选中英文单词
2. 等待浮动气泡出现（显示单词和释义）
3. 点击「收藏到 VocabWeb」按钮
4. 收藏成功后会显示提示消息

### 右键菜单

也可以通过右键菜单收藏：

1. 选中英文单词
2. 右键点击选中的文本
3. 选择「收藏到 VocabWeb」

### 查看统计

点击工具栏图标打开弹窗，可以查看：

- 今日收藏数量
- 总收藏数量
- 最近收藏的 5 个单词

## 开发调试

### 修改代码后重新加载

1. 修改代码文件
2. 访问 `chrome://extensions/`
3. 找到 VocabWeb 扩展
4. 点击刷新按钮 🔄

### 查看日志

- **Background Script 日志**：在扩展管理页面点击「Service Worker」查看
- **Content Script 日志**：在网页上按 F12 打开开发者工具查看
- **Popup 日志**：右键点击扩展图标 → 检查弹出内容

## 文件结构

```
extension/
├── manifest.json          # 扩展配置文件
├── background.js          # 后台服务 Worker
├── content.js             # 内容脚本（注入网页）
├── styles.css             # 内容脚本样式
├── popup.html             # 弹窗页面
├── popup.js               # 弹窗脚本
├── popup.css              # 弹窗样式
├── icons/                 # 图标文件
│   ├── icon16.png
│   ├── icon48.png
│   └── icon128.png
└── README.md              # 本文件
```

## 技术栈

- **Manifest V3**：Chrome 扩展最新版本
- **Vanilla JavaScript**：无框架依赖
- **Free Dictionary API**：免费的英语词典 API
- **Chrome Storage API**：本地存储配置

## 注意事项

- 扩展需要联网才能获取单词释义
- API Token 存储在本地，不会上传到任何服务器
- 浮动气泡使用高优先级样式，不会影响网页原有样式

## 许可证

MIT License
