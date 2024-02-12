import React, { useEffect } from "react";

import ReactFlow, { useNodesState, useEdgesState, Node, Edge } from "reactflow";
import Dagre from "@dagrejs/dagre";

import { Pipeline, PLtask } from "./bramble_types.ts";

// https://codesandbox.io/p/sandbox/romantic-bas-z2v5wm?file=%2FApp.js%3A63%2C51&utm_medium=sandpack
const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
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

const ExecutionView = (props: { pipeline: Pipeline }): React.ReactNode => {
    const pl: Pipeline = props.pipeline;
    const [nodes, setNodes, onNodesChange] = useNodesState(
        pl.spec.tasks ? generateNodes(pl.spec.tasks) : ([] as Node[])
    );
    const [edges, setEdges, onEdgesChange] = useEdgesState(
        pl.spec.tasks ? generateEdges(pl.spec.tasks) : ([] as Edge[])
    );

    useEffect(() => {
        const layouted = getLayoutedElements(nodes, edges);

        setNodes([...layouted.nodes]);
        setEdges([...layouted.edges]);
    });

    return (
        <>
            <h2 className="">Execution: {pl.metadata.name}</h2>
            <div
                className=""
                style={{
                    width: "100vw",
                    height: "50vh",
                    border: "black 3px",
                }}
            >
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    fitView={true}
                    onNodesChange={onNodesChange}
                    onEdgesChange={onEdgesChange}
                ></ReactFlow>
            </div>
        </>
    );
};

const generateNodes = (tasks: PLtask[]): Node[] => {
    return tasks.map((task: PLtask, i: number): Node => {
        return {
            id: task.name,
            position: { x: 0, y: i * 200 },
            data: {
                label: (
                    <div>
                        Name: {task.name}
                        <br></br>
                        Image: {task.spec.image}
                        <br></br>
                        Command: {task.spec.command.join(" ")}
                    </div>
                ),
            },
        };
    });
};

const generateEdges = (tasks: PLtask[]): Edge[] => {
    return tasks
        .map((task: PLtask): Edge[] => {
            return task.spec.dependencies
                ? task.spec.dependencies.map((dep, i): Edge => {
                      return {
                          id: task.name + i.toString(),
                          target: task.name,
                          source: dep,
                          type: "output",
                      };
                  })
                : ([] as Edge[]);
        })
        .flat();
};

export { ExecutionView };
