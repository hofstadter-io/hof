const data = [{
	name: "gen",
	title: "Code Generation",
	about: "CUE + text/template = _",
	items: [{
		name: "agnostic",
		about: "technology agnostic",
	},{
		name: "",
		about: "",
	}],
},{
	name: "datamodel",
	title: "Data Layer",
	about: "Evolvable data models and schemas",
	items: [{
		name: "one",
		about: "data model and constraints",
	},{

	}],
},{
	name: "flow",
	title: "Task Engine",
	about: "Define DAGs for data transforms, api calls, and arbitrary commands",
	items: [{
		name: "etl",
		about: "Build ETL pipelines",
	},{
		name: "",
		about: "",
	}],
},{
	name: "create",
	title: "Creators",
	about: "One line quick starts for your users",
	items: [{
		name: "",
		about: "",
	},{
		name: "",
		about: "",
	}],
}]

export default function C() {
  return (
		<>
			{ data.map( d => {
			return (
				<div id={d.name} className="flex w-full flex-col md:flex-row mb-8">
					<div className="border basis-full md:basis-1/2 p-2">
						<h1>{d.title}</h1>
						<p>{d.about}</p>
					</div>
					<div className="border basis-full md:basis-1/2 p-2">
						{ d.items.map( i => {
							return (
								<div id={i.name}>
								<p>{i.about}</p>
								</div>
							)
						})}
					</div>
				</div>
			)
			})}
		</>
	)
}

//export default function C() {
//  console.log(data)
//  return (
//    <div className="grid grid-cols-1 md:grid-cols-2">
//      { data.map( d => {
//      <>
//        <div className="">
//          <h3>{d.title}</h3>
//          <p>{d.about}</p>
//        </div>
//        <div>
//          { d.items.map( i => {
//            <p>{i.about}</p>
//          })}
//        </div>
//      </>
//      })}
//    </div>
//  )
//}
