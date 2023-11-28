import { createSignal } from "solid-js";
import "./App.css";
import * as bramble_types from "./bramble_types";

function App() {
  const [pipelines, setData] = createSignal<bramble_types.Pipeline[]>(
    new Array<bramble_types.Pipeline>(),
  );
  const fetchData = async () => {
    try {
      await fetch("http://localhost:5555/")
        .then((response) => response.json())
        .then((jsonData) =>
          setData(
            jsonData.map(
              (pipeline: bramble_types.Pipeline) =>
                new bramble_types.Pipeline(pipeline.metadata, pipeline.spec),
            ),
          ),
        );
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      {pipelines().map((pl?: bramble_types.Pipeline) => (
        <div class="pipeline">
          {
            <>
              <h2>Name: {pl?.metadata.name}</h2>
              <h3>Namespace: {pl?.metadata.namespace}</h3>
              <h2>Tasks</h2>
              <ul>
                {pl?.spec.tasks?.map((task: bramble_types.PLtask) => (
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
