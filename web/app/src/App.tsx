// import {useState} from 'react';
// import {Button} from 'primereact/button';
import { Drawer } from "@etxe/ui";

// function App() {
//     const [count,setCount] = useState(0);
//
//     return (
//         <div className="text-center">
//             <Button label="Click" icon="pi pi-plus" onClick={e => setCount(count + 1)}></Button>
//             <div className="text-2xl text-900 mt-3">{count}</div>
//         </div>
//     );
// }
//
// export default App;

export default function App() {
  return (
    <div>
      <Drawer width="200">bar</Drawer>
    </div>
  );
  // const [show, setShow] = useState();
  // return (
  //   <div className="App">
  //     <Button label="Click" icon="pi pi-plus" onClick={() => setShow((prevState) => !prevState)}>
  //       Click me
  //     </Button>
  //   </div>
  // );
}
