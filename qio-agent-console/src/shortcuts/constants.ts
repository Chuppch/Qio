export const WorkflowClipboardDataID = 'flowgram-workflow-clipboard-data';

export enum FlowCommandId {
  COPY = 'COPY', // 复制
  PASTE = 'PASTE', // 粘贴
  CUT = 'CUT', // 剪切
  GROUP = 'GROUP', // 组
  UNGROUP = 'UNGROUP', // 取消组
  COLLAPSE = 'COLLAPSE', // 折叠
  EXPAND = 'EXPAND', // 展开
  DELETE = 'DELETE', // 删除
  ZOOM_IN = 'ZOOM_IN', // 放大
  ZOOM_OUT = 'ZOOM_OUT', // 缩小
  RESET_ZOOM = 'RESET_ZOOM', // 重置缩放
  SELECT_ALL = 'SELECT_ALL', // 全选
  CANCEL_SELECT = 'CANCEL_SELECT', // 取消选中
}
