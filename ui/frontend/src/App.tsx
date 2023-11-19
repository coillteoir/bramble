import { createSignal } from 'solid-js'
import './App.css'
import { Pipeline } from './bramble_types'

function App() {
    const [pipelines, setData] = createSignal<[]>([])
    const fetchData = async () => {
        try {
            await fetch('http://localhost:5555/')
                 .then(response => response.json())
                 .then(jsonData => setData(jsonData))
        } catch (error) {
            console.error(error)
        }
    }
    setData([{
        metadata:{
            name: "test",
            namespace: "bramble",
        }, 
        spec:{
            tasks:[
                { 
                    name: "task01",
                    spec: {
                        image: "ubuntu",
                        command: ["sh", "-c", "hello world"]
                    },
                    dependencies: ["goodbye", "cheesy"]
                }
            ]
        }
    }])

  return (
    <>
    {
        pipelines().map(pl => 
            <div class="pipeline">
            {
                <>
                <h2>Name: {pl.metadata.name}</h2>
                <h3>Namespace: {pl.metadata.namespace}</h3>
                <h2>Tasks</h2>
                <ul>
                    {pl.spec.tasks.map(task => 
                    <div class="task">
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
