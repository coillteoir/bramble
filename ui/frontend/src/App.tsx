import { createSignal } from 'solid-js'
import './App.css'
import * as bramble_types from './bramble_types'

function App() {
    const [pipelines, setData] = createSignal<[bramble_types.Pipeline]>([])
    const fetchData = async () => {
        try {
            await fetch('http://localhost:5555/')
                 .then(response => response.json())
                 .then(jsonData => setData(jsonData))
        } catch (error) {
            console.error(error)
        }
    }
    setData([
    new bramble_types.Pipeline({name: "test", namespace: "bramble"},
        {
            tasks: [
                new bramble_types.PLtask("task01", 
                new bramble_types.TaskSpec("ubuntu", ["sh", "-c", "echo hello world"]),
                ["goodbye", "cheesy"])]
        }
    )])

  return (
    <>
    {
        pipelines().map(pl => 
            <div class="pipeline container bg-info">
            {
                <>
                <h2>Name: {pl.metadata.name}</h2>
                <h3>Namespace: {pl.metadata.namespace}</h3>
                <h2>Tasks</h2>
                <ul>
                    {pl.spec.tasks.map(task => 
                    <div class="task bg-primary">
                    <h2>{task.name}</h2>
                    <h3>Dependencies</h3>
                    <ul>{task.dependencies.map(dep => <li>{dep}</li>)}</ul>
                    </div>)}
                </ul>
                </>
            }
            </div>
        )
    }
    <button onclick={fetchData}>Get pipelines</button>
    </>
  )
}

export default App
