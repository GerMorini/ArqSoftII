import '../styles/Contact.css'
import gymPortada from '../assets/login/gimnasio1.jpeg'

const Contact = () => {
    const contacts = [
        {
            id: 1,
            name: "Germán Morini",
            email: "2302130@ucc.edu.ar",
            role: "API de usuario",
            initials: "GM",
            photo: "https://avatars.githubusercontent.com/u/97033560?v=4"
        },
        {
            id: 2,
            name: "Juan Fernandez Favarón",
            email: "2307329@ucc.edu.ar",
            role: "API de Actividades",
            initials: "FF",
            photo: "https://avatars.githubusercontent.com/u/163864272?v=4"
        },
        {
            id: 3,
            name: "Andrés García Villanueva",
            email: "2320413@ucc.edu.ar",
            role: "Frontend",
            initials: "AGV",
            photo: "https://avatars.githubusercontent.com/u/142047551?v=4"
        },
        {
            id: 4,
            name: "Pedro Giussano",
            email: "2102407@ucc.edu.ar",
            role: "API de búsqueda",
            initials: "PG",
            photo: "https://avatars.githubusercontent.com/u/163883332?v=4"
        }
    ];

    return (
        <div className="contact-container">
            <img
                className="contact-bg-img"
                src={gymPortada}
                alt="Fondo gimnasio"
            />

            <div className="contact-content">
                <div className="contact-header">
                    <h1 className="contact-title">
                        Nuestro <span className="gradient-text">Equipo</span>
                    </h1>
                    <p className="contact-subtitle">
                        Conoce a los profesionales que hicieron este proyecto realidad
                    </p>
                </div>

                <div className="contacts-list">
                    {contacts.map((contact, index) => (
                        <div
                            key={contact.id}
                            className={`contact-card ${index % 2 === 0 ? 'photo-left' : 'photo-right'}`}
                        >
                            <div className="contact-photo">
                                <div className="photo-placeholder">
                                    {contact.photo ? (
                                        <img
                                            src={contact.photo}
                                            alt={contact.name}
                                            className="contact-photo-img"
                                        />
                                    ) : (
                                        contact.initials
                                    )}
                                </div>
                            </div>
                            <div className="contact-details">
                                <h3 className="contact-name">{contact.name}</h3>
                                <p className="contact-role">{contact.role}</p>
                                <a href={`mailto:${contact.email}`} className="contact-email">
                                    📧 {contact.email}
                                </a>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Contact;
