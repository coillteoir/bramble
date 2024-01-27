import { createSignal, For } from "solid-js";
import { PipelineView } from "./PipelineView.tsx";
import { pipelinesBrambleDev } from "./bramble-types";

const App = () => {
    const [pipelines, setData] = createSignal<
        pipelinesBrambleDev.v1alpha1.Pipeline[]
    >(new Array<pipelinesBrambleDev.v1alpha1.Pipeline>());
    const fetchData = async (ns: string) => {
        try {
            await fetch("http://localhost:5555/pipelines/" + ns)
                .then((response) => response.json())
                .then((jsonData) => {
                    console.log("BACKEND RESPONSE:", jsonData[0]);
                    if (jsonData) {
                        setData(
                            jsonData.map(
                                (
                                    pipeline: any
                                ) => {
                                    const pl = new pipelinesBrambleDev.v1alpha1.Pipeline(
                                        pipeline
                                    )
                                    console.log("PIPELINE FROM CONSTRUCTOR:", pl)
                                    return pl
                                }
                            )
                        );
                    }
                });
        } catch (error) {
            console.error(error);
        }
    };
    const nsinput: any = <input class="" type="text" value="default" />;
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

                <button class="btn" onClick={() => fetchData(nsinput?.value)}>
                    Get pipelines
                </button>

                {pipelines()[0] && (
                    <h1>{pipelines()[0]?.metadata?.namespace}</h1>
                )}
                <For each={pipelines()}>
                    {(pl: pipelinesBrambleDev.v1alpha1.Pipeline) =>
                        pl && <PipelineView pipeline={pl} />
                    }
                </For>
            </div>
        </>
    );
};

export default App;
