import { createSignal, For } from "solid-js";
import "./App.css";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";

function App() {
    const [pipelines, setData] = createSignal<Pipeline[]>(
        new Array<Pipeline>()
    );
    const fetchData = async (ns: string) => {
        try {
            await fetch("http://localhost:5555/pipelines/" + ns)
                .then((response) => response.json())
                .then((jsonData) => {
                    setData(
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
    const nsinput: any = (
        <input class="rounded bg-black" type="text" value="default" />
    );
    return (
        <>
            <h1 class="absolute text-center inset-x-0 top-0 font-bold top">
                Bramble
            </h1>
            <label class="rounded bg-black">Namespace:</label>
            {nsinput}

            <button
                class="text-gray-400 hover:text-white"
                onClick={() => fetchData(nsinput?.value)}
            >
                Get pipelines
            </button>

            <h1>{pipelines()[0]?.metadata.namespace}</h1>
            <For each={pipelines()}>
                {(pl: Pipeline) => pl && <PipelineView pipeline={pl} />}
            </For>
        </>
    );
}

export default App;
