import { createSignal } from "solid-js";
import "./App.css";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline } from "./bramble_types";

function App() {
  const [pipelines, setData] = createSignal<Pipeline[]>(new Array<Pipeline>());
  const fetchData = async () => {
    try {
      await fetch("http://localhost:5555/pipelines")
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
  setInterval(fetchData, 10000)
  fetchData()
  return (
    <>
      <h1>{pipelines()[0]?.metadata.namespace}</h1>
      {pipelines().map((pl?: Pipeline) => (
              <PipelineView pipeline={pl} />
      ))}
      <button onclick={fetchData}>Get pipelines</button>
    </>
  );
}

export default App;
