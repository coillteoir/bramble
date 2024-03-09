import { PipelineView } from "./PipelineView.tsx";
import { pipelinesBrambleDev } from "./bramble-types";
import "reactflow/dist/style.css";
import React, { useState, useEffect, useRef } from "react";
import "./index.css";

const App = (): React.ReactNode => {
    const [namespace, setNamespace] = useState<string>("default");
    const [focusedPipeline, setFocusedPipeline] = useState<number>(0);
    const [pipelines, setPipelines] = useState<
        pipelinesBrambleDev.v1alpha1.Pipeline[]
    >(new Array<pipelinesBrambleDev.v1alpha1.Pipeline>());

    const inputRef = useRef<HTMLInputElement>(null);

    const fetchNamespacedPipelines = async () => {
        console.log("Fetching pipelines in:", namespace);
        try {
            await fetch("http://localhost:5555/pipelines/" + namespace)
                .then((response) => response.json())
                .then((jsonData) => {
                    console.log(jsonData)
                    setPipelines(
                        jsonData.map(
                            (pipeline: pipelinesBrambleDev.v1alpha1.Pipeline) =>
                                new pipelinesBrambleDev.v1alpha1.Pipeline(
                                    pipeline
                                )
                        )
                    );
                });
        } catch (error) {
            console.error(error);
        }
    };

    const fetchNamespacedPods = async () => {
        console.log("Fetching pods in:", namespace);
        try {
            await fetch("http://localhost:5555/pods/" + namespace)
                .then((response) => response.json())
                .then((jsonData) => {
                    console.log(jsonData);
                });
        } catch (error) {
            console.error(error);
        }
    };

    //setInterval(fetchNamespacedPods, 10000);
    //setInterval(fetchNamespacedPipelines, 10000);

    useEffect(() => {
        fetchNamespacedPods();
        fetchNamespacedPipelines();
    });

    return (
        <>
            <header className="">
                <h1 className="">Bramble</h1>
            </header>
            <div className="">
                <>
                    <input
                        className=""
                        type="text"
                        placeholder="Namespace"
                        ref={inputRef}
                    />

                    <button
                        className="btn btn-primary"
                        onClick={() => {
                            const ns: string | undefined =
                                inputRef.current?.value;
                            console.log(ns);
                            if (ns !== undefined) {
                                setNamespace(ns);
                            }
                        }}
                    >
                        Get pipelines
                    </button>
                    <div className="">
                        {pipelines.length !== 0 && (
                            <h1>Pipelines in the {namespace} namespace</h1>
                        )}
                        <ul
                            className=""
                        >
                            {pipelines.map(
                                (
                                    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline,
                                    index: number
                                ) =>
                                    pipeline && (
                                        <li
                                            key={index}
                                            className=""
                                            onClick={() => {
                                                setFocusedPipeline(index);
                                                console.log(
                                                    pipelines[focusedPipeline]
                                                );
                                            }}
                                        >
                                            {pipeline.metadata?.name}
                                        </li>
                                    )
                            )}
                        </ul>
                    </div>
                </>
                {pipelines.length !== 0 && pipelines[focusedPipeline] && (
                    <PipelineView pipeline={pipelines[focusedPipeline]} />
                )}
            </div>
        </>
    );
};

export default App;
