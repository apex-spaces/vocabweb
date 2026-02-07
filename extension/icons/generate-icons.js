// 生成简单的 PNG 图标占位文件
const fs = require('fs');
const path = require('path');

// 创建一个简单的 PNG 数据 URL（蓝色背景 + 白色 V 字母）
const createIconDataUrl = (size) => {
  // 使用 Canvas 或简单的 base64 编码的 PNG
  // 这里使用一个最小的 PNG 占位
  const canvas = `
    <svg width="${size}" height="${size}" xmlns="http://www.w3.org/2000/svg">
      <rect width="${size}" height="${size}" rx="${size/5}" fill="#4a9eff"/>
      <text x="${size/2}" y="${size*0.7}" font-family="Arial" font-size="${size*0.6}" font-weight="bold" fill="white" text-anchor="middle">V</text>
    </svg>
  `;
  return canvas;
};

// 输出 SVG 文件（浏览器可以直接使用）
const sizes = [16, 48, 128];
sizes.forEach(size => {
  const svg = createIconDataUrl(size);
  fs.writeFileSync(path.join(__dirname, `icon${size}.svg`), svg);
});

console.log('Icons created as SVG files');
