import stub from "@/init";
import React from "react";
import Sider from './Sider'
import Header from './Header'
import Content from './Content'

const Index: React.FC<any> = () => {
    const homeState = stub.ref.reactRedux.useSelector((state: any) => state.home)
    return (
        <div>
            {/*<stub.ref.intl.IntlProvider locale={"en"} messages={homeState.locale.message}>*/}
                <Sider/>
                <Header/>
                <Content/>
            {/*</stub.ref.intl.IntlProvider>*/}
        </div>

    );
}

export default Index