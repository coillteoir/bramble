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
