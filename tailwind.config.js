/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html,js}"],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
      colors: {
        "in-stock": "#78BE1F",
        "out-of-stock": "#F44336",
        order: "#FFC000",
        text: "#2F2F2F",
        subtext: "#C7C7C7",
        accent: "#FFC000",
        secondary: "#F7F8FA",
        delete: "#FF5F5F",
        border: "#D7DAE2",
        hover: "#5d5d5d",
        error: "#FF5F5F",
        "submit-disabled": "#D9D9D9",
        "delete-disabled": "#EE6460",
        "role-user": "#D9D9D9",
        "role-admin": "#FF5F5F",
        "role-artist": "#FFD700",
      },
      dropShadow: {
        glow: [
          "0 0px 20px rgba(255,255, 255, 0.35)",
          "0 0px 65px rgba(255, 255,255, 0.2)",
        ],
      },
    },
  },
  plugins: [],
};
