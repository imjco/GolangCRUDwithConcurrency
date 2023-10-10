


function editRow(){
  table.rows[rIndex].cells[1].innerHTML = document.getElementById("fname").value;
  table.rows[rIndex].cells[2].innerHTML = document.getElementById("lname").value;
  table.rows[rIndex].cells[3].innerHTML = document.getElementById("country").value;
  table.rows[rIndex].cells[4].innerHTML = document.getElementById("mNumber").value;
}

// Data Update Table Here
function editTableDisplay(){

  var table = document.getElementById("table_body"), rIndex;
  
for (var i = 1; i < table.rows.length; i++) {
 
  table.rows[i].onclick = function () {
    rIndex = this.rowIndex;
  
    // Get the values from the clicked row
    var id = this.cells[0].innerHTML;
    var title = this.cells[1].innerHTML;
    var author = this.cells[2].innerHTML;
    var isbn13 = this.cells[3].innerHTML;
    var isbn10 = this.cells[4].innerHTML;
    var publicationYear = this.cells[5].innerHTML;
    var publisher = this.cells[6].innerHTML;
    var edition = this.cells[7].innerHTML;
    var price = this.cells[8].innerHTML;

    // Populate the Edit Modal fields with the retrieved values
    document.getElementById("id").value = id;
    document.getElementById("title").value = title;
    document.getElementById("author").value = author;
    document.getElementById("isbn13").value = isbn13;
    document.getElementById("isbn10").value = isbn10;
    document.getElementById("publicationyear").value = publicationYear;
    document.getElementById("publisher").value = publisher;
    document.getElementById("edition").value = edition;
    document.getElementById("price").value = price;

   
  };
}

}

async function deleteBook(){
  let itemsTodelete = localStorage.getItem("todelete");
  let itemsTodeleteParse = JSON.parse(itemsTodelete);

  if (itemsTodeleteParse.isbn13 == '' || itemsTodeleteParse.isbn13 == null) {
     return
  }else{
    const response = await fetch(apiUrl + '/books/' + itemsTodeleteParse.isbn13, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      },
    })
  
    if(response.status === 200) {
      //location.reload();
      localStorage.removeItem('todelete')
      renderPage(currentPage);
      renderPagination();
      $('#deleteBookModal').modal('hide')
    }else{
      const textData = await response.text(); // Get response data as text
      alert(textData)
    }

  }





  //close modal
 

}
let bookDatabase = [];
let authorjsons = []



function readInfo(ID){
  $('#editBookModal').modal('show')

  const multiselect_u_a = document.getElementById("multiselect_u_a");
  const dropdownItems_u_a = document.getElementById("dropdownItems_u_a");
  const selectedItems_u_a = document.getElementById("selectedItems_u_a");
   
    multiselect_u_a.innerHTML = '';
		selectedItems_u_a.innerHTML = '';
    multiselect_u_p.innerHTML = '';
		selectedItems_u_p.innerHTML = '';
    
    // clear array selectedData_u_p
    selectedData_u_p.splice(0, selectedData_u_p.length);
    selectedData_u_a.splice(0, selectedData_u_a.length);



    document.getElementById("id_e").value = "";
    document.getElementById("title_e").value = "";
    document.getElementById("isbn10").value = "";
    document.getElementById("isbn13_e").value = "";
    document.getElementById("publicationyear_e").value = "";
    document.getElementById("edition_e").value = "";
    document.getElementById("price_e").value = "";
for (let i = 0; i < bookDatabase.length; i++) {

  if (bookDatabase[i].ID == ID) {
    document.querySelector('#id_e').value = bookDatabase[i].ID,
    document.querySelector('#title_e').value = bookDatabase[i].Title,
    document.querySelector("#isbn13_e").value = bookDatabase[i].ISBN13,
    document.querySelector("#isbn10_e").value = bookDatabase[i].ISBN10,
    document.querySelector("#publicationyear_e").value = bookDatabase[i].PublicationYear,
    document.querySelector("#edition_e").value = bookDatabase[i].Edition,
    document.querySelector("#price_e").value = bookDatabase[i].ListPrice
      
    for (let z = 0; z < bookDatabase[i].Authors.length; z++) {
      let authorname = bookDatabase[i].Authors[z].FirstName + " " + bookDatabase[i].Authors[z].LastName 
      appendToSelectedItems_UpdateAuthor(authorname, bookDatabase[i].Authors[z].ID)
    }
    for (let y = 0; y < bookDatabase[i].Publisher.length; y++) {
      appendToSelectedItems_UpdatePublisher(bookDatabase[i].Publisher[y].Name , bookDatabase[i].Publisher[y].ID)
    }
    break; // Exit the loop since we found a match
  }
}
}

let currentPage = 1;
let pageSize = 10;
let count = 10
let totalPages = Math.ceil(count / pageSize);


// Initial render
renderPage(currentPage);
renderPagination();

async function outputTableData(response) {
  
  const objectData = await response.json()
  count = objectData.BooksTotalCount;
  totalPages = Math.ceil(objectData.BooksTotalCount / pageSize)

  bookDatabase = objectData.Books

  let tableData=""
       objectData.Books.map((values)=>{
   
         let authorName = ""
         let Publisher = ""
         let authorjson= []
         for (let i = 0; i < values.Authors.length; i++) {
           authorName +=  values.Authors[i].FirstName + " " + values.Authors[i].MiddleName + " " + values.Authors[i].LastName
           if (i < values.Authors.length - 1) {
             authorName += ", "
           }
         }



         for (let i = 0; i < values.Authors.length; i++) {
          authorjson.push({
            "ID": values.Authors[i].ID,
            "FirstName": values.Authors[i].FirstName,
            "LastName": values.Authors[i].LastName,
            "MiddleName": values.Authors[i].MiddleName,
            "Books": null
          })
         }
       
         for (let i = 0; i < values.Publisher.length; i++) {
           Publisher +=  values.Publisher[i].Name
           if (i < values.Publisher.length - 1) {
             Publisher += ", "
           }
         }
      
       // add commas to price
        let convertPrice = values.ListPrice.toLocaleString();
       
        tableData+= `<tr> 
                 <td>${values.ID}</td>
                 <td>${values.Title}</td>
                 <td> ${authorName} </td>
                 <td>${values.ISBN13}</td>
                 <td>${values.ISBN10}</td>
                 <td>${values.PublicationYear}</td>
                 <td>${Publisher}</td>
                 <td>${values.Edition}</td>
                 <td>${convertPrice}</td>
                 <td></td>
                 <td>
                 <div class="text-center" ">
                     <div class="row">
                   
                      <button class="btn btn-warning btn-sm" onclick="readInfo('${values.ID}' )" data-bs-toggle="modal"><i class="material-icons" data-toggle="tooltip" title="Edit">&#xE254;</i></button>
                      
                     <button class="btn btn-danger btn-sm" onclick="showDeleteModal('${values.ID}', '${values.ISBN13}', '${values.Title}')"  data-bs-toggle="modal"><i class="material-icons" data-toggle="tooltip" title="Delete">&#xE872;</i></button>    
                     </div>
                 </div>
                 
                 </td>
                 </tr>`;
               });
               document.getElementById("table_body").innerHTML = tableData;
           renderPagination();

               let convertCount = count.toLocaleString();

               if (count === 0) {
                document.getElementById('showingcount').innerHTML = `            
                Showing <b>0</b> out of <b>${convertCount}</b> entries`
              }else{
                
                document.getElementById('showingcount').innerHTML = `Showing <b>${pageSize}</b> out of <b>${convertCount}</b> entries`
              }



              
}


async function renderPage(pageNum) {
    const response = await fetch( apiUrl + '/books?page='+pageNum +'&limit='+pageSize)
    if(response.status === 200) {
      
      outputTableData(response)
              
    } else{
      currentPage = 1;
      pageSize = 10;
      count = 0
      totalPages = Math.ceil(count / pageSize);

      renderPagination();
      document.getElementById("table_body").innerHTML = "";
    }

}


function showDeleteModal(id, isbn13, title){
  
  let data = {
    "id": id,
    "isbn13": isbn13,
    "title": title
  }
  
  localStorage.setItem("todelete", JSON.stringify(data));


  idTodelete = id;
  isbn13Todelete = isbn13;

  document.getElementById('deletetitle').innerHTML =`<h3>${title}</h3>`

   $('#deleteBookModal').modal('show')
}

function changepagesize(){
  pageSize = document.getElementById("pagesize").value;
  renderPage(currentPage);
  renderPagination();
}


function renderPagination() {
 
  // Render the pagination controls
 // var paginationElement = document.getElementById('myPagination');
  
 document.getElementById('myPagination').innerHTML=`
    <li class="page-item"><a class="page-link" href="#" onclick="prevPage()">Previous</a></li>
    ${renderPageNumbers()}
    <li class="page-item"><a class="page-link" href="#" onclick="nextPage()">Next</a></li>
  `;

  let convertCount = count.toLocaleString();
  if (count === 0) {
    document.getElementById('showingcount').innerHTML = `            
    Showing <b>0</b> out of <b>${convertCount}</b> entries`
  }else{
   
    document.getElementById('showingcount').innerHTML = `Showing <b>${pageSize}</b> out of <b>${convertCount}</b> entries`
  }
  



}
function renderPageNumbers() {

  // Render the page numbers for the current page and surrounding pages
  let startPage = Math.max(1, currentPage - 2);
  let endPage = Math.min(totalPages, currentPage + 2);

  let pageNumbers = '';
  for (let i = startPage; i <= endPage; i++) {
    pageNumbers += `<li class="page-item ${i === currentPage ? 'active' : ''}" data-pagenum="${i}" ><a class="page-link" href="#" onclick="goToPage(${i})">${i}</a></li>`;
  }

  return pageNumbers;
}
function goToPage(pageNum) {

  currentPage = pageNum;
  renderPage(currentPage);
  renderPagination();
}

function prevPage() {

  if (currentPage > 1) {
    currentPage--;
    renderPage(currentPage);
    renderPagination();
  }
}

function nextPage() {

  if (currentPage < totalPages) {
    currentPage++;
    renderPage(currentPage);
    renderPagination();
  }
}



async function addbook(){
  let title = document.getElementById("title").value;

  let isbn13 = document.getElementById("isbn13").value;
  let isbn10 = document.getElementById("isbn10").value;
  let publicationyear = document.getElementById("publicationyear").value;
  let edition = document.getElementById("edition").value;
  let price = document.getElementById("price").value;


  let selectedIds = selectedData.map((data) => data.id);
  let selectedIds_publisher = selectedData_p.map((data) => data.id);

  
let authordata = []
let publisherdata = []
for (let i = 0; i < selectedIds.length; i++) {
  authordata.push({
    "ID": parseInt(selectedIds[i]),
    "FirstName": "",
    "LastName": "",
    "MiddleName": "",
    "Books": null
  })
}

for (let a = 0; a < selectedIds_publisher.length; a++) {
  publisherdata.push({
    "ID": parseInt(selectedIds_publisher[a]),
    "Name": "",
    "Books": null
  })
}

  let data = {
    "Title": title,
    "Authors": authordata,
    "ISBN13": isbn13,
    "ISBN10": isbn10,
    "PublicationYear": parseInt(publicationyear),
    "Publisher": publisherdata,
    "Edition": edition,
    "ListPrice": parseInt(price)
  }

  const response = await fetch(apiUrl + '/books', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data),
  })

  if(response.status === 200) {
    location.reload();
   
    this.clearSelectedItems()
    renderPage(currentPage);
    renderPagination();
    $('#addBookModal').modal('hide')
    alert("Book Added Successfuly")
  }else{
    const textData = await response.text(); // Get response data as text
    alert(textData)
  }

}

async function search(){
  if(event.key === 'Enter') {
    var search = document.getElementById("myInput").value;
    // trim seach value
    var searchTrim = search.trim()
    // check if search is empty
    if (searchTrim !== "") {
      searchbyisbn(searchTrim)
    }else{
      currentPage = 1;
      renderPage(currentPage);
    }
}
}

async function searchbyisbn(value){
      const response = await fetch(apiUrl + '/books/' + value)
      if(response.status === 200) {
        outputTableData(response)
        const myInput = document.getElementById('myInput');
        myInput.value = '';
      }else{
        
      currentPage = 1;
      pageSize = 10;
      count = 0
      totalPages = Math.ceil(count / pageSize);

      renderPagination();
      document.getElementById("table_body").innerHTML = "";
      }
}



async function updatebook(){
  let id = document.getElementById("id_e").value;
  let title = document.getElementById("title_e").value;
  let isbn13 = document.getElementById("isbn13_e").value;
  let isbn10 = document.getElementById("isbn10_e").value;
  let publicationyear = document.getElementById("publicationyear_e").value;
  let edition = document.getElementById("edition_e").value;
  let price = document.getElementById("price_e").value;

  let selectedIds = selectedData_u_a.map((data) => data.id);
  let selectedIds_publisher = selectedData_u_p.map((data) => data.id);

let authordata = []
let publisherdata = []
for (let i = 0; i < selectedIds.length; i++) {
  authordata.push({
    "ID": parseInt(selectedIds[i]),
    "FirstName": "",
    "LastName": "",
    "MiddleName": "",
    "Books": null
  })
}

for (let a = 0; a < selectedIds_publisher.length; a++) {
  publisherdata.push({
    "ID": parseInt(selectedIds_publisher[a]),
    "Name": "",
    "Books": null
  })
}

  let data = {
    "ID": parseInt(id),
    "Title": title,
    "Authors": authordata,
    "ISBN13": isbn13,
    "ISBN10": isbn10,
    "PublicationYear": parseInt(publicationyear),
    "Publisher": publisherdata,
    "Edition": edition,
    "ListPrice": parseInt(price)
  }

  const response = await fetch(apiUrl + '/books/' + isbn13, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data),
  })

  if(response.status === 200) {
    location.reload();
   
   // clearSelectedItems()
    renderPage(currentPage);
    renderPagination();
    $('#addBookModal').modal('hide')
    alert("Book Update Successfuly")
  }else{
    const textData = await response.text(); // Get response data as text
    alert(textData)
  }

}

async function convertIsbn13(){
  var isbnInput = document.getElementById("AFrom");
  var isbnValue = isbnInput.value;

  // Remove any non-digit characters from the input
  isbnValue = isbnValue.replace(/\D/g, '');

  if (isbnValue.length === 13) {
     
     const response = await fetch( apiUrl + '/convertisbn13to10/' + isbnValue )
    if(response.status === 200) {
      const objectData = await response.json()
       document.getElementById("ATo").value = objectData.isbn10;
              
    } else{
      const textData = await response.text(); // Get response data as text
      alert(textData)
    }


  } else {
      alert("ISBN should contain exactly 13 digits.");
  }
}

async function convertIsbn10(){
  var isbnInput = document.getElementById("BFrom");
  var isbnValue = isbnInput.value;

  // Remove any non-digit characters from the input
  isbnValue = isbnValue.replace(/\D/g, '');

  if (isbnValue.length === 10) {
     
     const response = await fetch( apiUrl + '/convertisbn10to13/' + isbnValue )
    if(response.status === 200) {
      const objectData = await response.json()
       document.getElementById("BTo").value = objectData.isbn13;
              
    } else{
      const textData = await response.text(); // Get response data as text
      alert(textData)
    }


  } else {
      alert("ISBN should contain exactly 13 digits.");
  }
}