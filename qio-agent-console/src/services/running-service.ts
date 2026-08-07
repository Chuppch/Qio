import {
  injectable,
  inject,
  WorkflowDocument,
  Playground,
  delay,
  WorkflowLineEntity,
  WorkflowNodeEntity,
  WorkflowNodeLinesData,
} from '@flowgram.ai/free-layout-editor';
const RUNNING_INTERVAL = 1000;

@injectable()
export class RunningService {
  @inject(Playground) playground: Playground;

  @inject(WorkflowDocument) document: WorkflowDocument;

  private _runningNodes: WorkflowNodeEntity[] = [];

  async addRunningNode(node: WorkflowNodeEntity): Promise<void> {
    this._runningNodes.push(node); // 将节点加入到运行节点列表中
    node.renderData.node.classList.add('node-running'); // 访问当前节点的 DOM 元素，给 HTML 元素加一个 class 修改节点样式
    this.document.linesManager.forceUpdate(); // 强制让连线渲染器重新计算和刷新 —— 因为前面进行了 class 的调整
    await delay(RUNNING_INTERVAL);

    // Child Nodes
    await Promise.all(
      node.blocks.map((nextNode) => this.addRunningNode(nextNode))
    );
    // Sibling Nodes
    const nextNodes = node.getData(WorkflowNodeLinesData).outputNodes;
    // 并行执行子节点
    await Promise.all(
      nextNodes.map((nextNode) => this.addRunningNode(nextNode))
    );

  }

  async startRun(): Promise<void> {
    await this.addRunningNode(this.document.getNode('start_0')!);
    this._runningNodes.forEach((node) => {
      node.renderData.node.classList.remove('node-running');
    });
    this._runningNodes = [];
    this.document.linesManager.forceUpdate();
  }

  isFlowingLine(line: WorkflowLineEntity) {
    return this._runningNodes.some((node) =>
      node.getData(WorkflowNodeLinesData).outputLines.includes(line)
    );
  }
}
