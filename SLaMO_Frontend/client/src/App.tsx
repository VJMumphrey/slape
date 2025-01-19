import { BrowserRouter, Route, Routes } from "react-router-dom";
import index from "./pages/index.tsx";
import Dinosaur from "./pages/Dinosaur.tsx";
import { useState } from "react";
import "./App.css";
import Browser from "../../../../AppData/Local/deno/npm/registry.npmjs.org/debug/4.3.7/src/browser.js";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/:selectedDinosaur" element={Dinosaur />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
