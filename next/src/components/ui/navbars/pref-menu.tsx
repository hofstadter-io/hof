'use client';

import React from 'react';
import { FaCog } from "react-icons/fa";

export default function PrefMenu ({
		ThemeChanger,
		I18NChanger,
		CodeChanger,
	}:{
		ThemeChanger: React.ReactNode,
		I18NChanger: React.ReactNode,
		CodeChanger: React.ReactNode,
	}) {

	return (
		<div className="dropdown dropdown-end">
			<label tabIndex={0} className="btn btn-ghost btn-circle avatar">
				<FaCog className="w-1/2 h-1/2"/>
			</label>
			<ul tabIndex={0} className="mt-3 p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
				{ ThemeChanger ? <li><ThemeChanger /></li> : null }
			</ul>
		</div>
	)
}
