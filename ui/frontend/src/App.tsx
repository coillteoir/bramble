import { createSignal } from 'solid-js'
import './App.css'



function App() {
    const [data, setData] = createSignal<{name: string, image: string}[]>([{name: "goobert", image: "snoobert"}])
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
    {data().map(e => 
    <div class="pod">
        <p>{(`Name:${e.name} Image:${e.image}`)}</p>
    </div>)}
    <button onclick={fetchData}>Get data</button>
    </>
  )
}

export default App
