import { createSignal, For } from "solid-js";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";

const App = () => {
    const [namespace, setNamespace] = createSignal<string>("default");

    const [pipelines, setPipelines] = createSignal<Pipeline[]>(
        new Array<Pipeline>()
    );
    const fetchNamespacedPipelines = async (namespace: string) => {
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

    setInterval(fetchNamespacedPods, 1000, namespace());
    setInterval(fetchNamespacedPipelines, 1000, namespace());

    const nsinput: HTMLInputElement = (
        <input class="" type="text" />
    ) as HTMLInputElement;

    return (
        <>
            <header class="">
                <h1 class="">Bramble</h1>
                <ul class="">
                    <li class="">Pipelines</li>
                    <li class="">Tasks</li>
                </ul>
            </header>
            <div class="">
                <label class="">Namespace:</label>
                {nsinput}

                <button class="" onClick={() => setNamespace(nsinput.value)}>
                    Get pipelines
                </button>

                {pipelines().length !== 0 && (
                    <h1>
                        Pipelines in the {pipelines()[0]?.metadata.namespace}{" "}
                        namespace
                    </h1>
                )}
                <For each={pipelines()}>
                    {(pl: Pipeline) => pl && <PipelineView pipeline={pl} />}
                </For>
            </div>
        </>
    );
};

export default App;
