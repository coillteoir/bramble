import { createSignal } from 'solid-js'
import './App.css'



function App() {
    const [data, setData] = createSignal<string[]>(["Press button to get pods"])
    const fetchData = async () => {
        try {
            await fetch('http://localhost:5555/')
                 .then(response => response.json())
                 .then(jsonData => setData(jsonData))
        } catch (error) {
            console.error(error)
        }
    }

  return (
    <>
    <ul>
    {data().map(e => <li>{e}</li>)}
    </ul>
    <button onclick={fetchData}>Get data</button>
    </>
  )
}

export default App
