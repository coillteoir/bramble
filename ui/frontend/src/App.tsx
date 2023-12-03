import { createSignal } from "solid-js";
import "./App.css";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";

function App() {
  const [pipelines, setData] = createSignal<Pipeline[]>(new Array<Pipeline>());
  const fetchData = async (ns: string) => {
    try {
      await fetch("http://localhost:5555/pipelines/" + ns)
        .then((response) => response.json())
        .then((jsonData) => {
          setData(
            jsonData.map(
              (pipeline: Pipeline) =>
                new Pipeline(pipeline.metadata, pipeline.spec),
            ),
          );
        });
    } catch (error) {
      console.error(error);
    }
  };
  //setInterval(fetchData, 10000);
  //fetchData();
  const nsinput: any = <input type="text" value="default" />;
  return (
    <>
      <label>Namespace:</label>
      {nsinput}
      <button onclick={() => fetchData(nsinput.value)}>Get pipelines</button>
      <h1>{pipelines()[0]?.metadata.namespace}</h1>
      {pipelines().map((pl?: Pipeline) => (
        <PipelineView pipeline={pl} />
      ))}
    </>
  );
}

export default App;
