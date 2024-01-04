import { createSignal, For } from "solid-js";
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
    const nsinput: any = <input class="" type="text" value="default" />;
    return (
        <div class="">
            <h1 class="">Bramble</h1>
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
    );
}

export default App;
