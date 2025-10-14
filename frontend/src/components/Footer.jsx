import React from 'react'; 
import '../styles/Footer.css'; 

const Footer = () => {
    return (
        <footer className="footer-container">
            <p className="footer-text">
                &copy; {new Date().getFullYear()} Copyright GymPro
            </p>
        </footer>
    );
}

export default Footer;