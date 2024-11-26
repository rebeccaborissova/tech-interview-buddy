"use client";

import Link from "next/link";
import styles from "./Navbar.module.css"; 

const Navbar = () => {
  return (
    <nav className={styles.navbar}>
      <div className={styles.logo}>
        <Link href="/">CodeConnect</Link>
      </div>
      <ul className={styles.navLinks}>
        <li>
          <Link href="/login">Login</Link>
        </li>
        <li>
          <Link href="/sign-up">Sign Up</Link>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;
