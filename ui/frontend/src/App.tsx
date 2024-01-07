import { createSignal, For } from "solid-js";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";

const App = () => {
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
    const nsinput: any = <input class="" type="text" value="default" />;
    return (
        <>
            <header class="w-full bg-green-800 px-4">
                <h1 class="text-white font-sans text-2xl">Bramble</h1>
                <ul class="flex text-white">
                    <li class="inline-flex hover:bg-green-700 pr-4">
                        Pipelines
                    </li>
                    <li class="inline-flex hover:bg-green-700 px-4">Tasks</li>
                </ul>
            </header>
            <div class="px-4 bg-green-800">
                <label class="">Namespace:</label>
                {nsinput}

                <button class="" onClick={() => fetchData(nsinput?.value)}>
                    Get pipelines
                </button>

                <h1>{pipelines()[0]?.metadata.namespace}</h1>
                <For each={pipelines()}>
                    {(pl: Pipeline) => pl && <PipelineView pipeline={pl} />}
                </For>
            </div>
        </>
    );
};

export default App;
