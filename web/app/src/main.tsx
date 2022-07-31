// import React from 'react'

import 'primereact/resources/themes/mdc-light-indigo/theme.css';
import 'primereact/resources/primereact.min.css';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

// import {Metric} from "web-vitals/src/types";
// import reportWebVitals from './reportWebVitals';

import { createRoot } from 'react-dom/client'
import { BrowserRouter } from "react-router-dom"
import { App } from './App'

const root = createRoot(document.getElementById('root')!);
root.render(
  <BrowserRouter>
    <App overlayColor="white" width={600} />
  </BrowserRouter>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals();
