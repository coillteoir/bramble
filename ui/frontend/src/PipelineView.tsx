import React, { useEffect } from "react";

import ReactFlow, { useNodesState, useEdgesState, Node, Edge } from "reactflow";
import Dagre from "@dagrejs/dagre";

import { Pipeline, PLtask } from "./bramble_types.ts";

import {
    getLayoutedElements,
    generateNodes,
    generateEdges,
} from "./PipelineGraph.tsx";

// https://codesandbox.io/p/sandbox/romantic-bas-z2v5wm?file=%2FApp.js%3A63%2C51&utm_medium=sandpack
const PipelineView = (props: { pipeline: Pipeline }): React.ReactNode => {
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
            <h2 className="">Pipeline: {pl.metadata.name}</h2>
            <div
                className=""
                style={{
                    width: "100vw",
                    height: "100vh",
                    border: "black 3px",
                }}
            >
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    onNodesChange={onNodesChange}
                    onEdgesChange={onEdgesChange}
                    fitView
                    defaultEdgeOptions={{
                        animated: true,
                        style: {
                            "animation-direction": "reverse",
                            stroke: "green",
                        } as any,
                    }}
                />
            </div>
        </>
    );
};

export { PipelineView };
