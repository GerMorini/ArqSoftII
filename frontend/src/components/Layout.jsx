import Footer from "./Footer";
import Header from "./Header";
import { Outlet } from "react-router-dom";
import PageTransition from "./PageTransition";

const Layout = () => {
    return (
        <>
            <Header />
            <main className="main-content with-layout">
                <PageTransition>
                    <Outlet />
                </PageTransition>
            </main>
            <Footer />
        </>
    )
}

export default Layout;
