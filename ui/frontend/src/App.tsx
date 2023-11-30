import { createSignal } from "solid-js";
import "./App.css";
import { PipelineView } from "./PipelineView.tsx";
import { Pipeline, PLtask } from "./bramble_types";

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

  return (
    <>
      <h1>{pipelines()[0]?.metadata.namespace}</h1>
      {pipelines().map((pl?: Pipeline) => (
        <div class="pipeline">
          {
            <>
              <PipelineView pipeline={pl} />
              <h2>Name: {pl?.metadata.name}</h2>
              <h2>Tasks</h2>
              <ul>
                {pl?.spec.tasks?.map((task: PLtask) => (
                  <div class="task">
                    <h2>{task.name}</h2>
                    {task.dependencies && (
                      <>
                        <h3>Dependencies</h3>
                        <ul>
                          {task?.dependencies?.map((dep: string) => (
                            <li>{dep}</li>
                          ))}
                        </ul>
                      </>
                    )}
                  </div>
                ))}
              </ul>
            </>
          }
        </div>
      ))}
      <button onclick={fetchData}>Get pipelines</button>
    </>
  );
}

export default App;
