// import stub from "@/init";
import React from 'react';
import {BrowserRouter} from "react-router-dom";
import Home from "./layout/Home";

const App = () => {
    return (
        // <stub.ref.reactRedux.Provider store={stub.store}>
            <BrowserRouter basename={process.env.REACT_APP_BASE_URL}>
                <Home/>
            </BrowserRouter>
        // </stub.ref.reactRedux.Provider>
    );
}

export default App;
