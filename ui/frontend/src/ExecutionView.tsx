import React from "react";

import ReactFlow, { Node, Edge } from "reactflow";
import { Pod } from "kubernetes-models/v1";
import { pipelinesBrambleDev } from "./bramble-types";

import { getLayoutedElements } from "./Layout.tsx";
import { generateNodes, generateEdges } from "./ExecutionGraph.tsx";

const ExecutionView = (props: {
    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline;
    execution: pipelinesBrambleDev.v1alpha1.Execution;
    namespace: string;
}): React.ReactNode => {
    const [pods, setPods] = React.useState<Array<Pod>>(new Array<Pod>());

    const fetchNamespacedPods = async (namespace: string) => {
        console.log("Fetching pods in:", namespace);
        try {
            const baseUrl: string = "http://localhost:5555/";
            await fetch(baseUrl + "pods" + "/" + namespace)
                .then((response) => response.json())
                .then((jsonData) => {
                    setPods(jsonData);
                });
        } catch (error) {
            console.error(error);
        }
    };

    const pl: pipelinesBrambleDev.v1alpha1.Pipeline = props.pipeline;
    const layouted = getLayoutedElements(
        pl.spec?.tasks
            ? generateNodes(pl.spec?.tasks, pods, props.execution)
            : ([] as Node[]),
        pl.spec?.tasks ? generateEdges(pl.spec?.tasks) : ([] as Edge[])
    );

    const nodes = [...layouted.nodes];
    const edges = [...layouted.edges];

    React.useEffect(() => {
        fetchNamespacedPods(props.namespace);
    });

    return (
        <>
            <h2 className="">Pipeline: {pl.metadata?.name}</h2>
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
