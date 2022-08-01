import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import stub from "@/init";

const fetch: any = createAsyncThunk(
    'option/fetch',
    async (args: any) => {
        const res: any = await stub.api.post('option/query', {keys: args.keys})
        return res.data.data
    },
    {
        condition: (args, {getState, extra}) => {
            const {option} = getState() as any
            return args.force ?? stub.ref.lodash.differenceWith(args.keys, Object.keys(option.data), stub.ref.lodash.isEqual).length > 0
        }
    })

const slice = createSlice({
    name: "option",
    initialState: {
        data: {}
    },
    reducers: {},
    extraReducers: {
        [fetch.fulfilled]: (state: any, action: any) => {
            state.data = {...state.data, ...action.payload}
        }
    }
})
export default {
    slice,
    action: {...slice.actions, fetch}
}