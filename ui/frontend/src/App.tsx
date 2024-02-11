import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";
import "reactflow/dist/style.css";
import { useState, useEffect, useRef } from "react";

const App = () => {
    const [namespace, setNamespace] = useState<string>("default");
    const [pipelines, setPipelines] = useState<Pipeline[]>(
        new Array<Pipeline>()
    );

    const inputRef = useRef<HTMLInputElement>(null)

    const fetchNamespacedPipelines = async (namespace: string) => {
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

    const fetchNamespacedPods = async (namespace: string) => {
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

    //setInterval(fetchNamespacedPods, 10000, namespace);
    //setInterval(fetchNamespacedPipelines, 10000, namespace);

    useEffect(() => {
        fetchNamespacedPods(namespace);
        fetchNamespacedPipelines(namespace);
    });

    return (
        <>
            <header className="">
                <h1 className="">Bramble</h1>
                <ul className="">
                    <li className="">Pipelines</li>
                </ul>
            </header>
            <div className="">
                <label className="">Namespace:</label>
                <input className="" type="text" ref={inputRef} />

                <button
                    className=""
                    onClick={() => {
                        const ns: string | undefined = inputRef.current?.value;
                        console.log(ns);
                        if (ns !== undefined) {
                            setNamespace(ns);
                        }
                    }}
                >
                    Get pipelines
                </button>

                {pipelines.length !== 0 && (
                    <h1>
                        Pipelines in the {pipelines[0]?.metadata.namespace}{" "}
                        namespace
                    </h1>
                )}
                {pipelines.map(
                    (pl: Pipeline) => pl && <PipelineView pipeline={pl} />
                )}
            </div>
        </>
    );
};

export default App;
