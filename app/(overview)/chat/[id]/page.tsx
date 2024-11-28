type PageParams = {
  params: {id: string}
}

export default function Page({ params }: PageParams) {
  return<h1>{params.id}</h1>
}