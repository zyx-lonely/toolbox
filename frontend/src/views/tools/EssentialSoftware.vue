<template>
  <div>
    <n-h2>系统必备软件</n-h2>
    <n-p>精选常用系统必备软件，点击名称即可跳转官网下载。</n-p>

    <n-card>
      <n-tabs type="line" animated>
        <!-- 浏览器 -->
        <n-tab-pane name="browser" tab="🌐 浏览器">
          <n-data-table
            :columns="columns"
            :data="browserSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 办公软件 -->
        <n-tab-pane name="office" tab="📝 办公软件">
          <n-data-table
            :columns="columns"
            :data="officeSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 系统工具 -->
        <n-tab-pane name="tool" tab="🔧 系统工具">
          <n-data-table
            :columns="columns"
            :data="toolSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 安全软件 -->
        <n-tab-pane name="security" tab="🛡️ 安全软件">
          <n-data-table
            :columns="columns"
            :data="securitySoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 开发工具 -->
        <n-tab-pane name="dev" tab="💻 开发工具">
          <n-data-table
            :columns="columns"
            :data="devSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 设计创意 -->
        <n-tab-pane name="design" tab="🎨 设计创意">
          <n-data-table
            :columns="columns"
            :data="designSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 通讯社交 -->
        <n-tab-pane name="chat" tab="💬 通讯社交">
          <n-data-table
            :columns="columns"
            :data="chatSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 效率工具 -->
        <n-tab-pane name="efficiency" tab="⚡ 效率工具">
          <n-data-table
            :columns="columns"
            :data="efficiencySoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 媒体娱乐 -->
        <n-tab-pane name="media" tab="🎵 媒体娱乐">
          <n-data-table
            :columns="columns"
            :data="mediaSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 远程工具 -->
        <n-tab-pane name="remote" tab="🌐 远程工具">
          <n-data-table
            :columns="columns"
            :data="remoteSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 下载工具 -->
        <n-tab-pane name="download" tab="⬇️ 下载工具">
          <n-data-table
            :columns="columns"
            :data="downloadSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 游戏平台 -->
        <n-tab-pane name="game" tab="🎮 游戏平台">
          <n-data-table
            :columns="columns"
            :data="gameSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 虚拟机/容器 -->
        <n-tab-pane name="vm" tab="🖥️ 虚拟机">
          <n-data-table
            :columns="columns"
            :data="vmSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 输入法 -->
        <n-tab-pane name="input" tab="✏️ 输入法">
          <n-data-table
            :columns="columns"
            :data="inputSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 云存储 -->
        <n-tab-pane name="cloud" tab="☁️ 云存储">
          <n-data-table
            :columns="columns"
            :data="cloudSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 系统美化 -->
        <n-tab-pane name="beautify" tab="🎨 系统美化">
          <n-data-table
            :columns="columns"
            :data="beautifySoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>

        <!-- 硬件监控 -->
        <n-tab-pane name="hardware" tab="⚙️ 硬件监控">
          <n-data-table
            :columns="columns"
            :data="hardwareSoftware"
            :bordered="true"
            :single-line="false"
            size="small"
          />
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag } from 'naive-ui'

// 表格列定义（通用）
const columns = [
  {
    title: '软件名称',
    key: 'name',
    width: 180,
    render: (row: any) => h('a', {
      href: row.url,
      target: '_blank',
      style: 'color: #18a058; text-decoration: none; font-weight: bold; cursor: pointer;'
    }, { default: () => row.name })
  },
  { title: '描述', key: 'description', ellipsis: true },
  { title: '版本/大小', key: 'size', width: 100 },
  {
    title: '类型',
    key: 'free',
    width: 70,
    render: (row: any) => h(NTag, {
      type: row.free ? 'success' as const : 'warning' as const,
      size: 'small'
    }, { default: () => row.free ? '免费' : '付费' })
  }
]

// ==================== 浏览器 ====================
const browserSoftware = ref([
  { name: 'Google Chrome', description: 'Google 主流浏览器，插件生态丰富', size: '~100MB', url: 'https://www.google.com/chrome/', free: true },
  { name: 'Mozilla Firefox', description: '开源隐私浏览器，注重用户隐私保护', size: '~80MB', url: 'https://www.mozilla.org/firefox/', free: true },
  { name: 'Microsoft Edge', description: '微软浏览器，Win10/11 内置，Chromium 内核', size: '~150MB', url: 'https://www.microsoft.com/edge', free: true },
  { name: 'Opera / Opera GX', description: 'Opera 浏览器，GX 版专为游戏玩家设计', size: '~90MB', url: 'https://www.opera.com/', free: true },
  { name: 'Brave', description: '注重隐私的浏览器，内置广告拦截', size: '~110MB', url: 'https://brave.com/', free: true },
  { name: 'Vivaldi', description: '高度可定制的浏览器，面向高级用户', size: '~95MB', url: 'https://vivaldi.com/', free: true },
  { name: 'Arc (Windows)', description: '新一代浏览器，创新标签页管理方式', size: '~120MB', url: 'https://arc.net/', free: true },
  { name: 'Cent Browser', description: 'Chromium 内核浏览器，面向国内优化', size: '~100MB', url: 'https://centbrowser.com/', free: true }
])

// ==================== 办公软件 ====================
const officeSoftware = ref([
  { name: 'WPS Office', description: '金山办公套件，兼容 Office 格式，国产首选', size: '~200MB', url: 'https://www.wps.cn/', free: true },
  { name: 'LibreOffice', description: '开源办公套件，完全免费，功能强大', size: '~320MB', url: 'https://www.libreoffice.org/', free: true },
  { name: 'OnlyOffice', description: '现代风格的开源办公套件，界面美观', size: '~180MB', url: 'https://www.onlyoffice.com/', free: true },
  { name: 'Notepad++', description: '代码/文本编辑器，支持语法高亮和插件', size: '~4MB', url: 'https://notepad-plus-plus.org/', free: true },
  { name: 'VS Code', description: '微软轻量级代码编辑器，插件生态丰富', size: '~95MB', url: 'https://code.visualstudio.com/', free: true },
  { name: 'Typora', description: 'Markdown 编辑器，所见即所得', size: '~80MB', url: 'https://typora.io/', free: false },
  { name: 'Obsidian', description: '知识管理工具，双向链接笔记', size: '~100MB', url: 'https://obsidian.md/', free: true },
  { name: 'Sumatra PDF', description: '极速轻量 PDF 阅读器', size: '~10MB', url: 'https://www.sumatrapdfreader.org/', free: true },
  { name: 'PDF24 Creator', description: '免费 PDF 创建、合并、拆分工具', size: '~40MB', url: 'https://tools.pdf24.org/creator/', free: true },
  { name: 'Xmind', description: '思维导图软件', size: '~150MB', url: 'https://xmind.app/', free: false }
])

// ==================== 系统工具 ====================
const toolSoftware = ref([
  { name: '7-Zip', description: '开源高压缩比压缩解压工具', size: '~5MB', url: 'https://www.7-zip.org/', free: true },
  { name: 'Bandizip', description: '好用的压缩软件，界面友好', size: '~12MB', url: 'https://en.bandisoft.com/bandizip/', free: true },
  { name: 'Everything', description: '秒搜全盘文件的搜索神器', size: '~2MB', url: 'https://www.voidtools.com/', free: true },
  { name: 'Listary', description: '文件快速启动/搜索增强工具', size: '~8MB', url: 'https://www.listary.com/', free: true },
  { name: 'PowerToys', description: '微软官方 Windows 增强工具集（分屏、颜色取色等）', size: '~120MB', url: 'https://github.com/microsoft/PowerToys', free: true },
  { name: 'Geek Uninstaller', description: '极简强力卸载工具，体积小无残留', size: '~7MB', url: 'https://geekuninstaller.com/', free: true },
  { name: 'Bulk Crap Uninstaller', description: '批量卸载工具，功能强大', size: '~15MB', url: 'https://www.bcuninstaller.com/', free: true },
  { name: 'Revo Uninstaller', description: '专业卸载+注册表清理', size: '~30MB', url: 'https://www.revouninstaller.com/', free: false },
  { name: 'Dism++', description: 'Windows 系统精简/优化利器', size: '~15MB', url: 'https://www.chuyu.me/', free: true },
  { name: 'Process Lasso', description: 'CPU 进程优先级智能调整', size: '~10MB', url: 'https://bitsum.com/', free: false },
  { name: 'TreeSize Free', description: '磁盘空间可视化分析工具', size: '~8MB', url: 'https://www.jam-software.com/treesize_free/', free: true },
  { name: 'WinDirStat', description: '开源磁盘空间占用可视化分析', size: '~1MB', url: 'https://windirstat.net/', free: true },
  { name: 'Razer Cortex', description: '游戏加速/优化工具', size: '~200MB', url: 'https://www.razer.com/cortex', free: true },
  { name: 'QuickLook', description: '空格键预览文件（仿 macOS）', size: '~10MB', url: 'https://poooi.moe/quicklook/', free: true },
  { name: 'ClipboardFusion', description: '剪贴板历史管理/宏替换工具', size: '~25MB', url: 'https://www.binaryfortress.com/clipboardfusion/', free: true }
])

// ==================== 安全软件 ====================
const securitySoftware = ref([
  { name: '火绒安全', description: '轻量不打扰的杀毒+安全防护，国产推荐', size: '~50MB', url: 'https://www.huorong.cn/', free: true },
  { name: 'Malwarebytes', description: '专业恶意软件扫描清理工具', size: '~280MB', url: 'https://www.malwarebytes.com/', free: true },
  { name: 'HitmanPro.Alert', description: '行为分析型反勒索/防漏洞工具', size: '~15MB', url: 'https://www.hitmanpro.com/', free: false },
  { name: 'GlassWire', description: '网络防火墙监控可视化', size: '~20MB', url: 'https://www.glasswire.com/', free: false },
  { name: 'VeraCrypt', description: '开源磁盘加密工具（TrueCrypt 继任者）', size: '~35MB', url: 'https://veracrypt.fr/', free: true },
  { name: 'Bitwarden', description: '开源跨平台密码管理器', size: '~10MB', url: 'https://bitwarden.com/', free: true },
  { name: 'KeePassXC', description: '本地密码数据库管理器', size: '~18MB', url: 'https://keepassxc.org/', free: true }
])

// ==================== 开发工具 ====================
const devSoftware = ref([
  { name: 'Git', description: '分布式版本控制系统', size: '~50MB', url: 'https://git-scm.com/', free: true },
  { name: 'GitHub Desktop', description: 'Git 图形化客户端', size: '~90MB', url: 'https://desktop.github.com/', free: true },
  { name: 'Node.js', description: 'JavaScript 运行时环境', size: '~30MB', url: 'https://nodejs.org/', free: true },
  { name: 'Python', description: '编程语言及环境', size: '~50MB', url: 'https://python.org/', free: true },
  { name: 'JDK (Oracle)', description: 'Java 开发工具包', size: '~160MB', url: 'https://www.oracle.com/java/technologies/downloads/', free: true },
  { name: 'Go', description: 'Google 编程语言及工具链', size: '~130MB', url: 'https://go.dev/dl/', free: true },
  { name: 'Rust', description: '系统编程语言及工具链', size: '~150MB', url: 'https://rustup.rs/', free: true },
  { name: 'Docker Desktop', description: '容器化开发环境', size: '~500MB', url: 'https://www.docker.com/products/docker-desktop/', free: true },
  { name: 'Postman', description: 'API 调试与测试工具', size: '~140MB', url: 'https://www.postman.com/', free: true },
  { name: 'DevToys', description: '开发者瑞士军刀工具箱（编码转换、JSON格式化等）', size: '~10MB', url: 'https://devtoys.app/', free: true },
  { name: 'Windows Terminal', description: '微软现代终端（PowerShell/CMD/WSL）', size: '~20MB', url: 'https://github.com/microsoft/terminal', free: true },
  { name: 'Oh My Posh', description: '漂亮的终端提示符主题', size: '~5MB', url: 'https://ohmyposh.dev/', free: true }
])

// ==================== 设计创意 ====================
const designSoftware = ref([
  { name: 'Figma', description: '在线协作 UI/UX 设计工具', size: '~0(网页)', url: 'https://www.figma.com/', free: true },
  { name: 'GIMP', description: '开源图像编辑软件（Photoshop 替代品）', size: '~230MB', url: 'https://www.gimp.org/', free: true },
  { name: 'Krita', description: '开源数字绘画软件', size: '~150MB', url: 'https://krita.org/', free: true },
  { name: 'Inkscape', description: '开源矢量图形编辑器（Illustrator 替代品）', size: '~100MB', url: 'https://inkscape.org/', free: true },
  { name: 'DaVinci Resolve', description: '专业视频剪辑/调色软件（免费版功能已很强）', size: '~800MB', url: 'https://www.blackmagicdesign.com/products/davinciresolve', free: true },
  { name: 'OBS Studio', description: '开源录屏/直播推流软件', size: '~120MB', url: 'https://obsproject.com/', free: true },
  { name: 'ShareX', description: '开源录屏截图/GIF/上传一体化工具', size: '~20MB', url: 'https://getsharex.com/', free: true },
  { name: 'Snipaste', description: '高效截图/贴图工具', size: '~15MB', url: 'https://zh.snipaste.com/', free: true },
  { name: 'ScreenToGif', description: '轻量 GIF 录制/编辑工具', size: '~5MB', url: 'https://screentogif.com/', free: true },
  { name: 'Paint.NET', description: '轻量易用的图片编辑器', size: '~10MB', url: 'https://www.getpaint.net/', free: true }
])

// ==================== 通讯社交 ====================
const chatSoftware = ref([
  { name: '微信 (WeChat)', description: '腾讯即时通讯软件', size: '~250MB', url: 'https://weixin.qq.com/', free: true },
  { name: 'QQ', description: '腾讯 QQ 客户端', size: '~200MB', url: 'https://im.qq.com/', free: true },
  { name: '钉钉 (DingTalk)', description: '阿里企业通讯协作平台', size: '~300MB', url: 'https://www.dingtalk.com/', free: true },
  { name: '飞书 (Lark)', description: '字节跳动企业协作平台', size: '~250MB', url: 'https://www.feishu.cn/', free: true },
  { name: '企业微信', description: '微信企业版', size: '~300MB', url: 'https://work.weixin.qq.com/', free: true },
  { name: 'Discord', description: '游戏社区语音聊天平台', size: '~130MB', url: 'https://discord.com/', free: true },
  { name: 'Telegram', description: '加密即时通讯软件', size: '~60MB', url: 'https://telegram.org/', free: true },
  { name: 'Slack', description: '团队协作沟通工具', size: '~140MB', url: 'https://slack.com/', free: true }
])

// ==================== 效率工具 ====================
const efficiencySoftware = ref([
  { name: 'AutoHotkey', description: 'Windows 自动化脚本/热键工具', size: '~5MB', url: 'https://www.autohotkey.com/', free: true },
  { name: 'Power Automate', description: '微软 RPA 自动化工具（Win11 内置）', size: '~20MB', url: 'https://powerautomate.microsoft.com/', free: true },
  { name: 'Ditto', description: '开源剪贴板历史管理工具', size: '~3MB', url: 'http://ditto-cp.sourceforge.net/', free: true },
  { name: 'Rainmeter', description: 'Windows 桌面自定义美化工具', size: '~5MB', url: 'https://rainmeter.net/', free: true },
  { name: 'f.lux', description: '屏幕色温自动调节护眼软件', size: '~2MB', url: 'https://justgetflux.com/', free: true },
  { name: 'EarTrumpet', description: 'Windows 音量混合器增强工具', size: '~5MB', url: 'https://eartrumpet.app/', free: true },
  { name: 'TranslucentTB', description: '任务栏透明效果工具', size: '~1MB', url: 'https://github.com/TranslucentTB/TranslucentTB', free: true },
  { name: 'Files (App)', description: '现代风格的文件资源管理器替代品', size: '~80MB', url: 'https://files.community/', free: true },
  { name: 'CopyQ', description: '高级剪贴板管理器，支持编辑/脚本', size: '~15MB', url: 'https://hluk.github.io/CopyQ/', free: true },
  { name: 'Keyboard Manager (PowerToys)', description: '键盘快捷键重映射（含在 PowerToys 中）', size: '-', url: 'https://github.com/microsoft/PowerToys', free: true }
])

// ==================== 媒体娱乐 ====================
const mediaSoftware = ref([
  { name: 'VLC Media Player', description: '开源万能播放器', size: '~42MB', url: 'https://www.videolan.org/vlc/', free: true },
  { name: 'PotPlayer', description: '功能强大的本地播放器', size: '~35MB', url: 'https://potplayer.daum.net/', free: true },
  { name: 'MPC-BE / MPC-HC', description: '经典轻量播放器', size: '~15MB', url: 'https://mpc-be.github.io/', free: true },
  { name: 'foobar2000', description: '高品质音乐播放器', size: '~6MB', url: 'https://www.foobar2000.org/', free: true },
  { name: 'mpv Player', description: '轻量命令行播放器，可高度定制', size: '~20MB', url: 'https://mpv.io/', free: true },
  { name: 'AIMP', description: '精美音乐播放器', size: '~15MB', url: 'https://aimp.ru/', free: true },
  { name: 'Audacity', description: '开源音频录制编辑工具', size: '~35MB', url: 'https://www.audacityteam.org/', free: true },
  { name: 'handbrake', description: '开源视频转码压缩工具', size: '~30MB', url: 'https://handbrake.fr/', free: true },
  { name: 'FFmpeg', description: '命令行音视频处理工具集', size: '~50MB', url: 'https://ffmpeg.org/', free: true },
  { name: 'yt-dlp', description: '命令行视频下载工具（YouTube 等）', size: '~10MB', url: 'https://github.com/yt-dlp/yt-dlp', free: true },
  { name: 'Stacher', description: 'yt-dlp 的 GUI 封装，图形化下载视频', size: '~60MB', url: 'https://github.com/Stacher/stacher', free: true },
  { name: 'ImageMagick', description: '命令行图片批量处理工具', size: '~20MB', url: 'https://imagemagick.org/', free: true },
  { name: 'Caesium', description: '图片批量压缩工具', size: '~8MB', url: 'https://saerasoft.com/caesium/', free: true }
])

// ==================== 远程工具 ====================
const remoteSoftware = ref([
  { name: 'ToDesk', description: '国产远程桌面，免费高速', size: '~25MB', url: 'https://www.todesk.com/', free: true },
  { name: '向日葵远程控制', description: '贝锐旗下远程桌面工具', size: '~30MB', url: 'https://sunlogin.oray.com/', free: true },
  { name: 'AnyDesk', description: '轻量远程桌面工具', size: '~5MB', url: 'https://anydesk.com/', free: true },
  { name: 'TeamViewer', description: '老牌远程桌面工具', size: '~60MB', url: 'https://www.teamviewer.com/', free: true },
  { name: 'RustDesk', description: '开源远程桌面（TeamViewer 替代品）', size: '~15MB', url: 'https://rustdesk.com/', free: true },
  { name: 'UltraVNC', description: '开源 VNC 远程桌面', size: '~3MB', url: 'https://uvnc.com/', free: true },
  { name: 'Chrome Remote Desktop', description: 'Google 远程桌面（浏览器插件）', size: '~2MB', url: 'https://remotedesktop.google.com/', free: true },
  { name: 'Microsoft Remote Desktop (RDP)', description: 'Windows 自带远程桌面客户端', size: '~0(内置)', url: 'https://apps.microsoft.com/store/detail/microsoft-remote-desktop/9WZDNCRFJ3P', free: true }
])

// ==================== 下载工具 ====================
const downloadSoftware = ref([
  { name: 'IDM (Internet Download Manager)', description: '老牌多线程下载加速工具', size: '~15MB', url: 'https://www.internetdownloadmanager.com/', free: false },
  { name: 'FDM (Free Download Manager)', description: '开源免费下载加速工具', size: '~50MB', url: 'https://www.freedownloadmanager.org/', free: true },
  { name: 'qBittorrent', description: '开源 BT/磁力下载工具', size: '~40MB', url: 'https://www.qbittorrent.org/', free: true },
  { name: 'Transmission', description: '轻量 BT 下载工具', size: '~15MB', url: 'https://transmissionbt.com/', free: true },
  { name: 'Aria2', description: '命令行多线程下载工具', size: '~2MB', url: 'https://aria2.github.io/', free: true },
  { name: 'Motrix', description: '开源下载工具（Aria2 图形界面）', size: '~80MB', url: 'https://motrix.app/', free: true },
  { name: 'EagleGet', description: '免费下载加速工具', size: '~15MB', url: 'https://www.eagleget.com/', free: true },
  { name: 'XDown', description: '磁力/种子/直链下载工具', size: '~5MB', url: 'https://xdown.org/', free: true }
])

// ==================== 游戏平台 ====================
const gameSoftware = ref([
  { name: 'Steam', description: 'Valve 游戏平台', size: '~500MB', url: 'https://store.steampowered.com/', free: true },
  { name: 'Epic Games Launcher', description: 'Epic 游戏平台（常送免费游戏）', size: '~200MB', url: 'https://www.epicgames.com/store/', free: true },
  { name: 'GOG Galaxy', description: 'GOG 游戏平台（无 DRM）', size: '~150MB', url: 'https://www.gog.com/galaxy', free: true },
  { name: 'Xbox (Win11 内置)', description: '微软游戏平台（Win11 内置）', size: '~0(内置)', url: 'https://www.xbox.com/', free: true },
  { name: 'Playnite', description: '游戏库统一管理工具（整合 Steam/Epic 等）', size: '~200MB', url: 'https://playnite.link/', free: true },
  { name: 'MSI Afterburner', description: '显卡超频/监控工具', size: '~50MB', url: 'https://www.msi.com/Landing/afterburner/', free: true },
  { name: 'Razer Cortex', description: '游戏加速/优化工具', size: '~200MB', url: 'https://www.razer.com/cortex', free: true }
])

// ==================== 虚拟机/容器 ====================
const vmSoftware = ref([
  { name: 'VMware Workstation Player', description: 'VMware 免费虚拟机', size: '~500MB', url: 'https://www.vmware.com/products/workstation-player.html', free: true },
  { name: 'Oracle VirtualBox', description: '开源虚拟机软件', size: '~220MB', url: 'https://www.virtualbox.org/', free: true },
  { name: 'Docker Desktop', description: '容器化开发环境', size: '~500MB', url: 'https://www.docker.com/products/docker-desktop/', free: true },
  { name: 'WSL2 (Windows 内置)', description: 'Windows 子系统 Linux', size: '~0(内置)', url: 'https://learn.microsoft.com/windows/wsl/', free: true },
  { name: 'QEMU + Virt-Manager', description: '开源虚拟化方案', size: '~300MB', url: 'https://www.qemu.org/', free: true }
])

// ==================== 输入法 ====================
const inputSoftware = ref([
  { name: '微软拼音/五笔 (Win10/11 内置)', description: 'Windows 自带输入法', size: '~0(内置)', url: 'https://support.microsoft.com/windows', free: true },
  { name: '搜狗输入法', description: '流行中文输入法', size: '~100MB', url: 'https://shurufa.sogou.com/', free: true },
  { name: '讯飞输入法', description: '语音输入+智能输入', size: '~150MB', url: 'https://www.ifo x.com/', free: true },
  { name: '手心输入法', description: '无广告纯净输入法', size: '~50MB', url: 'https://www.xinshuru.com/', free: true },
  { name: '小狼毫 (Rime)', description: '开源高度可定制输入法', size: '~80MB', url: 'https://rime.im/', free: true }
])

// ==================== 云存储/同步 ====================
const cloudSoftware = ref([
  { name: 'OneDrive', description: '微软云存储（Win10/11 内置）', size: '~0(内置)', url: 'https://onedrive.live.com/', free: true },
  { name: 'Google Drive', description: 'Google 云存储', size: '~150MB', url: 'https://www.google.com/drive/', free: true },
  { name: 'Dropbox', description: '老牌云存储同步工具', size: '~120MB', url: 'https://www.dropbox.com/', free: true },
  { name: '百度网盘', description: '国内云存储（需安装客户端）', size: '~200MB', url: 'https://pan.baidu.com/', free: true },
  { name: '阿里云盘', description: '阿里云存储（速度较快）', size: '~150MB', url: 'https://www.aliyundrive.com/', free: true },
  { name: 'Syncthing', description: '开源无服务器文件同步工具', size: '~30MB', url: 'https://syncthing.net/', free: true },
  { name: 'FreeFileSync', description: '开源文件夹对比同步工具', size: '~20MB', url: 'https://freefilesync.org/', free: true }
])

// ==================== 系统美化/增强 ====================
const beautifySoftware = ref([
  { name: 'Start11', description: '开始菜单/任务栏美化工具', size: '~100MB', url: 'https://www.stardock.com/products/start11/', free: false },
  { name: 'OpenShell', description: '开源开始菜单自定义工具', size: '~10MB', url: 'https://github.com/Open-Shell/Open-Shell-Menu', free: true },
  { name: 'TranslucentTB', description: '任务栏透明效果工具', size: '~1MB', url: 'https://github.com/TranslucentTB/TranslucentTB', free: true },
  { name: 'Rainmeter', description: '桌面小部件/皮肤美化工具', size: '~5MB', url: 'https://rainmeter.net/', free: true },
  { name: 'Wallpaper Engine', description: '动态壁纸工具（Steam）', size: '~500MB', url: 'https://store.steampowered.com/app/431960/', free: false },
  { name: 'Lively Wallpaper', description: '开源动态壁纸工具', size: '~100MB', url: 'https://github.com/rocksdanister/lively/', free: true },
  { name: 'ElectricSheep', description: '开源动态壁纸（休眠羊群）', size: '~20MB', url: 'https://electricsheep.org/', free: true }
])

// ==================== 硬件监控/测试 ====================
const hardwareSoftware = ref([
  { name: 'CPU-Z', description: 'CPU/主板/内存信息检测工具', size: '~5MB', url: 'https://www.cpuid.com/softwares/cpu-z.html', free: true },
  { name: 'GPU-Z', description: '显卡信息检测工具', size: '~10MB', url: 'https://www.techpowerup.com/gpuz/', free: true },
  { name: 'HWiNFO64', description: '全面的硬件信息/传感器监控工具', size: '~20MB', url: 'https://www.hwinfo.com/', free: true },
  { name: 'CrystalDiskInfo', description: '硬盘健康状态检测工具', size: '~10MB', url: 'https://crystalmark.info/en/software/crystaldiskinfo/', free: true },
  { name: 'CrystalDiskMark', description: '硬盘速度基准测试工具', size: '~5MB', url: 'https://crystalmark.info/en/software/crystaldiskmark/', free: true },
  { name: 'FurMark', description: '显卡压力测试/烤机工具', size: '~20MB', url: 'https://www.gpumark.com/', free: true },
  { name: 'Prime95', description: 'CPU 压力测试/烤机工具', size: '~10MB', url: 'https://www.mersenne.org/download/', free: true },
  { name: 'MemTest86', description: '内存稳定性测试工具', size: '~10MB', url: 'https://www.memtest86.com/', free: true }
])
</script>
