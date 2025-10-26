import Footer from "./Footer";
import Header from "./Header";
import { Outlet } from "react-router-dom";
import PageTransition from "./PageTransition";
import { useTokenValidation } from "../hooks/useTokenValidation";

const Layout = ({ setAlertDialog }) => {
    // Verificar la expiración del token cada vez que se carga una página
    useTokenValidation(setAlertDialog);

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
