<script lang="ts">
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import Pen from '@lucide/svelte/icons/pen';
	import Plus from '@lucide/svelte/icons/plus';
	import Export from '@lucide/svelte/icons/folder-down';
	import Print from '@lucide/svelte/icons/printer';
	import ContactIcon from '@lucide/svelte/icons/contact';
	import axios from 'axios';

	let contacts: { [key: string]: string | number }[] = $state([
		{
			ID: 0,
			FirstName: 'Max',
			LastName: 'Mustermann',
			Street: 'Musterallee',
			HouseNumber: '600',
			ZipCode: '12345',
			City: 'Berlin',
			Partner: '-',
			Children: '-',
			Email: 'maxmustermann@gmail.com'
		}
	]);

	axios.get('http://localhost:8080/contacts').then((response) => {
		contacts = response.data;
	}).catch((error) => {
		console.log(error);
	});

	let searchTerm = $state('');
	let showEdit = $state(false);
	let showNew = $state(false);
	let editContact = $state(0);
	let newContact: { [key: string]: string } = $state({});

	function valueInDictValues(dict: { [key: string]: string | number }, searchValue: string): boolean {
		for (const key in dict) {
			if (typeof dict[key] === 'string' && dict[key].includes(searchValue)) {
				return true;
			}
		}
		return false;
	}
</script>

<div class="fixed top-0 w-full z-50 bg-white shadow-md">
	<AppBar>
		{#snippet lead()}
			<ContactIcon size={32} />
		{/snippet}
		{#snippet trail()}
			<button
				class="bg-gray-500 hover:bg-gray-600 text-white p-2 rounded-md"
				aria-label="Print"
				onclick={() => {
					let tableHtml = "<table>";
					tableHtml += "<thead><tr><th>Vorname</th><th>Nachname</th><th>Straße</th><th>Hausnummer</th><th>Postleitzahl</th><th>Ort</th><th>Partner</th><th>Kinder</th><th>E-Mail</th></tr></thead>";
					tableHtml += "<tbody>";
					for (const item of contacts) {
  				  tableHtml += `<tr><td>${item.FirstName ? item.FirstName : "-"}</td><td>${item.LastName ? item.LastName : "-"}</td><td>${item.Street ? item.Street : "-"}</td><td>${item.HouseNumber ? item.HouseNumber : "-"}</td><td>${item.ZipCode ? item.ZipCode : "-"}</td><td>${item.City ? item.City : "-"}</td><td>${item.Partner ? item.Partner : "-"}</td><td>${item.Children ? item.Children : "-"}</td><td>${item.Email ? item.Email : "-"}</td></tr>`;
  				}
					tableHtml += "</tbody></table>";
  				const printWindow = window.open("", "_blank");
  				if (printWindow) {
  				  const doc = printWindow.document;
  				  doc.head.innerHTML = `
  				    <title>Kontakte</title>
  				    <style>
  				      table {
  				        width: 100%;
  				        border-collapse: collapse;
  				        font-family: Arial, sans-serif;
  				      }
  				      th, td {
  				        border: 1px solid #ddd;
  				        padding: 8px;
  				        text-align: left;
  				      }
  				      th {
  				        background-color: #f0f0f0;
  				        font-weight: bold;
  				      }
  				      tbody tr:nth-child(even) {
  				        background-color: #f9f9f9;
  				      }
  				      @media print {
  				        body {
  				          font-size: 12pt;
  				        }
  				        table {
  				          width: 100% !important;
  				        }
  				      }
  				    </style>
  				  `;
  				  doc.body.innerHTML = tableHtml;
  				  doc.close();
  				  printWindow.print();
  				} else {
  				  alert("Konnte Druckfenster nicht öffnen.");
  				}
					}}>
				<Print size={16} />
			</button>
			<button
				class="bg-gray-500 hover:bg-gray-600 text-white p-2 rounded-md"
				aria-label="exportToCSV"
				onclick={() => {
					axios.get('http://localhost:8080/csvExport', {responseType: "blob"}).then((response) => {
						console.log(response);
						const blob = response.data;
						const url = URL.createObjectURL(blob);
						const a = document.createElement('a');
						a.href = url;
						a.download = 'kontakte.csv';
						document.body.appendChild(a);
						a.click();
						window.URL.revokeObjectURL(url);
						document.body.removeChild(a);
					}).catch((error) => {
						console.log(error);
					})
					}}>
				<Export size={16} />
			</button>
			<button
				class="bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-md"
				aria-label="newButton"
				onclick={() => {
						showNew = true;
						newContact = {};
					}}>
				<Plus size={16} />
			</button>
			<div class="relative w-full max-w-md">
				<input
					type="text"
					placeholder="Search..."
					class="w-full pl-4 pr-12 py-2 border rounded-md focus:outline-none focus:ring focus:ring-blue-300"
					bind:value={searchTerm}
				/>
				<button class="absolute top-0 right-0 h-full px-4 text-gray-600 hover:text-gray-800"
								aria-label="searchButton">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
				</button>
			</div>
		{/snippet}
		<span class="text-3xl font-bold">Kontakte</span>
	</AppBar>
</div>

<div class="flex flex-wrap gap-4 justify-center items-center pt-[74px]">
	{#each contacts as contact, id (contact.ID)}
		{#if valueInDictValues(contact, searchTerm)}
			<div
				class="relative card preset-outlined-primary-500 w-full max-w-md p-4 m-3 text-center grid grid-cols-3 gap-4">
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Vorname</div>
					<div class="text-base">{contact.FirstName ? contact.FirstName : "-"}</div>
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Nachname</div>
					<div class="text-base">{contact.LastName ? contact.LastName : "-"}</div>
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Partner</div>
					<div class="text-base">{contact.Partner ? contact.Partner : "-"}</div>
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Straße</div>
					<div class="text-base">{contact.Street ? contact.Street : "-"}</div>
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Hausnummer</div>
					<div class="text-base">{contact.HouseNumber ? contact.HouseNumber : "-"}</div>
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Ort</div>
					<div class="text-base">{contact.City ? contact.City : "-"}</div>
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Postleitzahl</div>
					<div class="text-base">{contact.ZipCode ? contact.ZipCode : "-"}</div>
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">E-Mail</div>
					<div class="text-base">{contact.Email ? contact.Email : "-"}</div>
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Kinder</div>
					<div class="text-base">{contact.Children ? contact.Children : "-"}</div>
				</div>

				<button
					class="absolute top-2 right-2 opacity-0 hover:opacity-100 transition-opacity duration-300 bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-md z-10"
					aria-label="editButton"
					onclick={() => {
						showEdit = true;
						editContact = id;
						}
					}>
					<Pen size={16} />
				</button>
			</div>
		{/if}
	{/each}

	{#if showEdit}
		<div class="fixed inset-0 flex items-center justify-center z-10">
			<div
				class="relative card preset-outlined-primary-500 bg-white w-full max-w-md p-4 m-3 text-center grid grid-cols-3 gap-4">
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Vorname</div>
					<input class="input" type="text" placeholder="Vorname"
								 bind:value={contacts[editContact].FirstName} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Nachname</div>
					<input class="input" type="text" placeholder="Nachname"
								 bind:value={contacts[editContact].LastName} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Partner</div>
					<input class="input" type="text" placeholder="Partner" bind:value={contacts[editContact].Partner} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Straße</div>
					<input class="input" type="text" placeholder="Straße" bind:value={contacts[editContact].Street} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Hausnummer</div>
					<input class="input" type="text" placeholder="Hausnummer"
								 bind:value={contacts[editContact].HouseNumber} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Ort</div>
					<input class="input" type="text" placeholder="Ort" bind:value={contacts[editContact].City} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Postleitzahl</div>
					<input class="input" type="text" placeholder="Postleitzahl"
								 bind:value={contacts[editContact].ZipCode} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">E-Mail</div>
					<input class="input" type="text" placeholder="E-Mail" bind:value={contacts[editContact].Email} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Kinder</div>
					<input class="input" type="text" placeholder="Kinder" bind:value={contacts[editContact].Children} />
				</div>

				<button onclick={() => {
					showEdit = false;
					window.location.reload();}
					}
								class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded mr-2">
					Abbrechen
				</button>
				<button onclick={() => {
					showEdit = false;
					axios.delete("http://localhost:8080/deleteContact" , {data: JSON.stringify({ID: contacts[editContact].ID})})
					window.location.reload();
				}
					}
								class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded mr-2">
					Löschen
				</button>
				<button onclick={() => {
					showEdit = false;
					axios.put("http://localhost:8080/updateContact" , JSON.stringify(contacts[editContact]));
				}
					}
								class="bg-primary-500 hover:bg-primary-700 text-white font-bold py-2 px-4 rounded mr-2">
					Ändern
				</button>
			</div>
		</div>
	{/if}

	{#if showNew}
		<div class="fixed inset-0 flex items-center justify-center z-10">
			<div
				class="relative card preset-outlined-primary-500 bg-white w-full max-w-md p-4 m-3 text-center grid grid-cols-3 gap-4">
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Vorname</div>
					<input class="input" type="text" placeholder="Vorname" bind:value={newContact.FirstName} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Nachname</div>
					<input class="input" type="text" placeholder="Nachname" bind:value={newContact.LastName} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Partner</div>
					<input class="input" type="text" placeholder="Partner" bind:value={newContact.Partner} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Straße</div>
					<input class="input" type="text" placeholder="Straße" bind:value={newContact.Street} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Hausnummer</div>
					<input class="input" type="text" placeholder="Hausnummer" bind:value={newContact.HouseNumber} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">Ort</div>
					<input class="input" type="text" placeholder="Ort" bind:value={newContact.City} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Postleitzahl</div>
					<input class="input" type="text" placeholder="Postleitzahl" bind:value={newContact.ZipCode} />
				</div>
				<div class="col-span-2 relative">
					<div class="text-[0.6rem] text-gray-500">E-Mail</div>
					<input class="input" type="text" placeholder="E-Mail" bind:value={newContact.Email} />
				</div>
				<div class="relative">
					<div class="text-[0.6rem] text-gray-500">Kinder</div>
					<input class="input" type="text" placeholder="Kinder" bind:value={newContact.Children} />
				</div>

				<button onclick={() => showNew = false}
								class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded mr-2">
					Abbrechen
				</button>
				<div></div>
				<button onclick={() => {
					showNew = false;
					axios.post("http://localhost:8080/newContact", JSON.stringify($state.snapshot(newContact))).then(() => {
						window.location.reload();})
				}}
								class="bg-primary-500 hover:bg-primary-700 text-white font-bold py-2 px-4 rounded mr-2">
					Hinzufügen
				</button>
			</div>
		</div>
	{/if}
</div>
