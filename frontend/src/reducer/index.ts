import Home from './Home'
import Option from './Option'
import {configureStore} from "@reduxjs/toolkit";

const reducer = {
    reducer: {
        home: Home.slice.reducer,
        option: Option.slice.reducer
    },
    action: {
        home: Home.action,
        option: Option.action
    }
}
const store = configureStore({reducer: reducer.reducer})

export {
    store,
    reducer
}