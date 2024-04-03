import React from "react";

import ReactFlow, { Node, Edge } from "reactflow";
import { Pod } from "kubernetes-models/v1";
import { pipelinesBrambleDev } from "./bramble-types";

import { getLayoutedElements } from "./Layout.tsx";
import { generateNodes, generateEdges } from "./ExecutionGraph.tsx";

const ExecutionView = (props: {
    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline;
    execution: pipelinesBrambleDev.v1alpha1.Execution | undefined;
}): React.ReactNode => {
    const [pods, setPods] = React.useState<Array<Pod>>(new Array<Pod>());

    const fetchNamespacedPods = async (namespace: string | undefined) => {
        if (namespace === undefined) {
            console.log("undefined namespace");
            return;
        }

        console.log(`Fetching pods in ${namespace} ${Date()}`);
        try {
            await fetch(`http://localhost:5555/resources/pods/${namespace}`)
                .then((response) => response.json())
                .then((jsonData) => {
                    setPods(jsonData);
                });
        } catch (error) {
            console.error(error);
        }
    };

    const layouted = getLayoutedElements(
        props.pipeline.spec?.tasks
            ? generateNodes(props.pipeline.spec?.tasks, pods, props.execution)
            : ([] as Node[]),
        props.pipeline.spec?.tasks
            ? generateEdges(props.pipeline.spec?.tasks)
            : ([] as Edge[])
    );

    const nodes = [...layouted.nodes];
    const edges = [...layouted.edges];

    React.useEffect(() => {
        const interval = setInterval(
            () => fetchNamespacedPods(props.execution?.metadata?.namespace),
            1000
        );
        return () => {
            clearInterval(interval);
        };
    });

    return (
        <>
            <h2 className="">Pipeline: {props.pipeline.metadata?.name}</h2>
            {props.execution && (
                <h2 className="">
                    Execution: {props.execution.metadata?.name}
                </h2>
            )}
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
