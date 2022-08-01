import React from "react"
import stub from '@/init'

export interface MenuModel {
    name: string,
    path?: string
    icon?: React.ReactNode
    children?: Array<MenuModel>
}

const menus: MenuModel[] = [
    {
        name: "环境",
        path: "/environment",
        // icon: <stub.ref.icon.ai.ClusterOutlined/>,
    },
    {
        name: "集群",
        path: "/cluster",
        // icon: <stub.ref.icon.ai.ClusterOutlined/>,
    },
    {
        name: "租户",
        path: "/tenant",
        // icon: <stub.ref.icon.ai.ClusterOutlined/>,
    },
    {
        name: "命名空间",
        path: "/namespace",
        // icon: <stub.ref.icon.ai.ClusterOutlined/>,
    }
];
export default menus