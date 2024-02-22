import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";
import "reactflow/dist/style.css";
import React, { useState, useEffect, useRef } from "react";
import "./index.css";

const App = (): React.ReactNode => {
    const [namespace, setNamespace] = useState<string>("default");
    const [focusedPipeline, setFocusedPipeline] = useState<number>(0);
    const [pipelines, setPipelines] = useState<Pipeline[]>(
        new Array<Pipeline>()
    );

    const inputRef = useRef<HTMLInputElement>(null);

    const fetchNamespacedPipelines = async () => {
        console.log("Fetching pipelines in:", namespace);
        try {
            await fetch("http://localhost:5555/pipelines/" + namespace)
                .then((response) => response.json())
                .then((jsonData) => {
                    setPipelines(
                        jsonData.map(
                            (pipeline: Pipeline) =>
                                new Pipeline(pipeline.metadata, pipeline.spec)
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
    }, []);

    return (
        <>
            <header className="bg-green-400">
                <h1 className="text-3xl font-bold">Bramble</h1>
                <ul className="">
                    <li className="">Pipelines</li>
                </ul>
            </header>
            <div className="bg-green-200 w-screen h-screen">
                <>
                    <input
                        className=""
                        type="text"
                        placeholder="Namespace"
                        ref={inputRef}
                    />

                    <button
                        className="
                        text-white 
                        bg-green-800 
                        rounded"
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
                    <div className="divide-y divide-slate-400">
                        {pipelines.length !== 0 && (
                            <h1>Pipelines in the {namespace} namespace</h1>
                        )}
                        <ul
                            className="
                            divide-y 
                            divide-slate-400 
                            w-fit"
                        >
                            {pipelines.map(
                                (pipeline: Pipeline, index: number) =>
                                    pipeline && (
                                        <li
                                            key={index}
                                            className="
                                            hover:border-10 
                                            hover:border-black
                                            hover:border-solid"
                                            onClick={() => {
                                                setFocusedPipeline(index);
                                                console.log(
                                                    pipelines[focusedPipeline]
                                                );
                                            }}
                                        >
                                            {pipeline.metadata.name}
                                        </li>
                                    )
                            )}
                        </ul>
                    </div>
                </>
                {pipelines.length !== 0 && (
                    <PipelineView pipeline={pipelines[focusedPipeline]} />
                )}
            </div>
        </>
    );
};

export default App;
