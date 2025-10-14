import Footer from "./Footer";
import Header from "./Header";
import { Outlet } from "react-router-dom";

const Layout = () => {
    return (
        <> 
            <Header />
            <main className="main-content with-layout"> 
                <Outlet />
            </main>
            <Footer />
        </>
    )
}

export default Layout;
