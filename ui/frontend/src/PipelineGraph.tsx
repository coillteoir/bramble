import { PLtask } from "./bramble_types.ts";

import { Node, Edge } from "reactflow";
import Dagre from "@dagrejs/dagre";

// https://codesandbox.io/p/sandbox/romantic-bas-z2v5wm?file=%2FApp.js%3A63%2C51&utm_medium=sandpack
export const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
    const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
    g.setGraph({ rankdir: "TB" });

    edges.forEach((edge: Edge) => g.setEdge(edge.source, edge.target));

    // any type used because the typing of setNode seems to be incorrect.
    // It behaves as expected when entire node is passed in
    nodes.forEach((node: Node) => g.setNode(node.id, node as any));

    Dagre.layout(g);

    return {
        nodes: nodes.map((node: Node) => {
            const { x, y } = g.node(node.id);
            return { ...node, position: { x, y } };
        }),
        edges,
    };
};

export const generateNodes = (tasks: PLtask[]): Node[] => {
    return tasks.map((task: PLtask): Node => {
        return {
            id: task.name,
            width: 120,
            height: 50,
            position: { x: 0, y: 0 },
            data: {
                label: (
                    <div>
                        <p>{task.name}</p>
                    </div>
                ),
            },
        };
    });
};

export const generateEdges = (tasks: PLtask[]): Edge[] => {
    return tasks
        .map((task: PLtask): Edge[] => {
            return task.spec.dependencies
                ? task.spec.dependencies.map((dep: string, i: number): Edge => {
                      return {
                          id: task.name + i.toString(),
                          target: dep,
                          source: task.name,
                      };
                  })
                : ([] as Edge[]);
        })
        .flat();
};
