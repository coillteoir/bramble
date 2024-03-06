import React from "react";

import ReactFlow, { Node, Edge } from "reactflow";

import { Pipeline } from "./bramble_types.ts";

import {
    getLayoutedElements,
    generateNodes,
    generateEdges,
} from "./PipelineGraph.tsx";

// https://codesandbox.io/p/sandbox/romantic-bas-z2v5wm?file=%2FApp.js%3A63%2C51&utm_medium=sandpack
const ExecutionView = (props: { pipeline: Pipeline }): React.ReactNode => {
    const pl: Pipeline = props.pipeline;
    const layouted = getLayoutedElements(
        pl.spec.tasks ? generateNodes(pl.spec.tasks) : ([] as Node[]),
        pl.spec.tasks ? generateEdges(pl.spec.tasks) : ([] as Edge[])
    );

    const nodes = [...layouted.nodes];
    const edges = [...layouted.edges];

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
                    fitView
                    defaultEdgeOptions={{
                        animated: true,
                        style: {
                            animationDirection: "reverse",
                            stroke: "green",
                        },
                    }}
                />
            </div>
        </>
    );
};

export { ExecutionView };
