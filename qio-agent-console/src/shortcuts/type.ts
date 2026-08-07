import type { WorkflowJSON } from '@flowgram.ai/free-layout-editor';

import type { WorkflowClipboardDataID } from './constants';

/**
 * 剪贴板数据来源信息
 * 用于标记数据的来源环境，防止跨环境复制粘贴导致配置错误。
 * 不同环境的数据库配置不同，相同 ID 可能指向不同的资源。
 */
export interface WorkflowClipboardSource {
  /** 当前页面的域名+端口 (window.location.host)，如: localhost:3000 */
  host: string;
  // more: id?, workspaceId? etc.
}

/**
 * 节点边界信息
 * 用于粘贴时智能定位：计算偏移量并保持节点间相对位置。
 */
export interface WorkflowClipboardRect {
  x: number;
  y: number;
  width: number;
  height: number;
}

/**
 * 剪贴板数据结构*
 * 复制节点时保存到系统剪贴板的完整数据，粘贴时读取并验证。
 */
export interface WorkflowClipboardData {
  type: typeof WorkflowClipboardDataID;
  json: WorkflowJSON;
  source: WorkflowClipboardSource;
  bounds: WorkflowClipboardRect;
}
