import {createSlice} from "@reduxjs/toolkit";
import locale from '@/locale'

interface Locale {
    key: string
    message: any
}

const slice = createSlice({
    name: "home",
    initialState: {
        collapsed: false,
        locale: {
            key: "en_US",
            message: locale["en_US"]
        } as Locale,
        environment: {}
    },
    reducers: {
        collapse: (state) => {
            return {...state, ...{collapsed: !state.collapsed}};
        },
        switchLocale: (state, action) => {
            return {...state, locale: {...state.locale, ...action.payload}}
        },
        setEnvironment: (state, action) => {
            return {...state, environment: action.payload};
        }
    }
})
export default {
    slice,
    action: {...slice.actions}
}