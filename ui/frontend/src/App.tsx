import { PipelineView } from "./PipelineView.tsx";
import { ExecutionView } from "./ExecutionView.tsx";
import { PipelineList } from "./PipelineList.tsx";
import { pipelinesBrambleDev } from "./bramble-types";
import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import "reactflow/dist/style.css";
import React, { useState, useEffect, useRef } from "react";
import "./index.css";

const NamespaceSearch = (props: {
    inputRef: React.RefObject<HTMLInputElement>;
    setNamespace: React.Dispatch<React.SetStateAction<string>>;
}): React.ReactNode => (
    <div className="">
        {/*
        <div className="label">
            <span className="label-text">Enter Namespace</span>
        </div>
       */}

        <input
            className="input input-bordered input-primary w-3/5"
            type="text"
            ref={props.inputRef}
            placeholder="Enter Namespace"
        />

        <button
            className="btn btn-primary w-2/5"
            onClick={() => {
                const ns: string | undefined = props.inputRef.current?.value;
                console.log(ns);
                if (ns !== undefined && ns !== "") {
                    props.setNamespace(ns);
                }
            }}
        >
            Get pipelines
        </button>
    </div>
);

const View = (props: {
    pipelines: Array<Pipeline>;
    focusedPipeline: number;
    executions: Array<Execution>;
    focusedExecution: string | undefined;
    namespace: string;
}): React.ReactNode =>
    (!props.focusedExecution &&
        props.pipelines.length !== 0 &&
        props.pipelines[props.focusedPipeline] && (
            <PipelineView pipeline={props.pipelines[props.focusedPipeline]} />
        )) ||
    (props.focusedExecution && (
        <ExecutionView
            pipeline={props.pipelines[props.focusedPipeline]}
            execution={
                props.executions.filter(
                    (exe: Execution) =>
                        exe.metadata?.name === props.focusedExecution
                )[0]
            }
            namespace={props.namespace}
        />
    ));

const App = (): React.ReactNode => {
    const [namespace, setNamespace] = useState<string>("default");

    const [focusedPipeline, setFocusedPipeline] = useState<number>(0);
    const [focusedExecution, setFocusedExecution] = useState<
        string | undefined
    >("");

    const [pipelines, setPipelines] = useState<Pipeline[]>(
        new Array<Pipeline>()
    );

    const [executions, setExecutions] = useState<
        pipelinesBrambleDev.v1alpha1.Execution[]
    >(new Array<Execution>());

    const inputRef = useRef<HTMLInputElement>(null);

    const fetchNamespacedResources = async () => {
        const fetchRes = async (resource: string) => {
            console.log("Fetching", resource, "in:", namespace);
            try {
                const baseUrl: string = "http://localhost:5555/";
                await fetch(baseUrl + resource + "/" + namespace)
                    .then((response) => response.json())
                    .then((jsonData) => {
                        switch (resource) {
                            case "pipelines": {
                                setPipelines(jsonData);
                                break;
                            }
                            case "executions": {
                                setExecutions(jsonData);
                                break;
                            }
                        }
                    });
            } catch (error) {
                console.error(error);
            }
        };
        await fetchRes("pipelines");
        await fetchRes("executions");
    };

    useEffect(() => {
        fetchNamespacedResources();
    }, [namespace]);

    return (
        <>
            <React.StrictMode>
                <header className="">
                    <h1 className="text-3xl font-bold bg-slate-800">Bramble</h1>
                </header>
                <div className="w-1/3 lg:w-1/6">
                    <NamespaceSearch
                        inputRef={inputRef}
                        setNamespace={setNamespace}
                    />

                    <PipelineList
                        namespace={namespace}
                        pipelines={pipelines}
                        executions={executions}
                        focusedPipeline={focusedPipeline}
                        setFocusedPipeline={setFocusedPipeline}
                        setFocusedExecution={setFocusedExecution}
                    />
                </div>
                <View
                    pipelines={pipelines}
                    focusedPipeline={focusedPipeline}
                    executions={executions}
                    focusedExecution={focusedExecution}
                    namespace={namespace}
                />
            </React.StrictMode>
        </>
    );
};

export default App;
